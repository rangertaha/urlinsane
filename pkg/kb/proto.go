// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package kb

import (
	"fmt"
	"sort"
	"strconv"

	"google.golang.org/protobuf/proto"

	"github.com/rangertaha/urlinsane/pkg/kb/internal/kbpb"
)

// The dataset is stored as protocol buffers rather than as the JSON a Layout
// marshals to. The wire form is a third of the size, decodes without
// reflection over field names, and is embedded into the binary, so the cost of
// the whole catalogue is paid in bytes rather than in parse time. The JSON
// methods remain on Layout and Key for callers that want to read or export a
// layout; nothing in the package reads them back.

var (
	formToProto = map[Form]kbpb.Form{
		ANSI: kbpb.Form_FORM_ANSI,
		ISO:  kbpb.Form_FORM_ISO,
		JIS:  kbpb.Form_FORM_JIS,
	}

	formFromProto = map[kbpb.Form]Form{
		kbpb.Form_FORM_ANSI: ANSI,
		kbpb.Form_FORM_ISO:  ISO,
		kbpb.Form_FORM_JIS:  JIS,
	}
)

// Scan codes and virtual keys are the two things every key carries, so the
// dataset stores them as numbers rather than as the "10" and "VK_Q" spellings
// the source data uses. Both are numbers to begin with — a scan code is the
// byte the keyboard reports, and a virtual key is a Windows constant — so the
// strings were never more than a presentation of them, and they cost around a
// hundred kilobytes across the catalogue.
//
// The public API keeps the readable spelling. These four functions are the
// only place the two representations meet.

// encodeSC parses the hex spelling of a scan code.
func encodeSC(sc string) (uint32, error) {
	n, err := strconv.ParseUint(sc, 16, 32)
	if err != nil {
		return 0, fmt.Errorf("scan code %q is not hex", sc)
	}
	return uint32(n), nil
}

// decodeSC spells a scan code the way the source data does, in at least two
// uppercase hex digits.
func decodeSC(sc uint32) string { return fmt.Sprintf("%02X", sc) }

// encodeKLID parses the hex spelling of a keyboard layout identifier.
func encodeKLID(klid string) (uint32, error) {
	n, err := strconv.ParseUint(klid, 16, 32)
	if err != nil {
		return 0, fmt.Errorf("KLID %q is not hex", klid)
	}
	return uint32(n), nil
}

// decodeKLID spells a KLID the way Windows does, as eight lowercase hex
// digits. Lookups depend on this: the catalogue is keyed by the result.
func decodeKLID(klid uint32) string { return fmt.Sprintf("%08x", klid) }

// encodeVK looks up a virtual key by name. An unnamed key is stored as
// unspecified; an unrecognised one is an error, so that a virtual key the
// schema has never seen stops the build rather than being quietly dropped.
func encodeVK(vk string) (kbpb.VirtualKey, error) {
	if vk == "" {
		return kbpb.VirtualKey_VK_UNSPECIFIED, nil
	}
	n, ok := kbpb.VirtualKey_value[vk]
	if !ok {
		return 0, fmt.Errorf("unknown virtual key %q", vk)
	}
	return kbpb.VirtualKey(n), nil
}

// decodeVK names a virtual key.
func decodeVK(vk kbpb.VirtualKey) (string, error) {
	if vk == kbpb.VirtualKey_VK_UNSPECIFIED {
		return "", nil
	}
	name, ok := kbpb.VirtualKey_name[int32(vk)]
	if !ok {
		return "", fmt.Errorf("unknown virtual key %d", vk)
	}
	return name, nil
}

// encodeLocales converts the language identities a layout is installed under.
func encodeLocales(locales []Locale) ([]*kbpb.Locale, error) {
	out := make([]*kbpb.Locale, 0, len(locales))
	for _, loc := range locales {
		klid, err := encodeKLID(loc.KLID)
		if err != nil {
			return nil, err
		}
		out = append(out, &kbpb.Locale{Klid: klid, Tag: loc.Tag, Name: loc.Name})
	}
	return out, nil
}

// decodeLocales is the inverse of encodeLocales.
func decodeLocales(msgs []*kbpb.Locale) []Locale {
	out := make([]Locale, 0, len(msgs))
	for _, loc := range msgs {
		out = append(out, Locale{KLID: decodeKLID(loc.Klid), Tag: loc.Tag, Name: loc.Name})
	}
	return out
}

// Marshal encodes the layout in the dataset's storage format.
func (l *Layout) Marshal() ([]byte, error) {
	form, ok := formToProto[l.Form]
	if !ok {
		return nil, fmt.Errorf("kb: layout %q has unknown form %q", l.ID, l.Form)
	}

	locales, err := encodeLocales(l.Locales)
	if err != nil {
		return nil, fmt.Errorf("kb: layout %q: %w", l.ID, err)
	}

	msg := &kbpb.Layout{
		Id:      l.ID,
		Name:    l.Name,
		File:    l.File,
		Form:    form,
		Locales: locales,
		Keys:    make([]*kbpb.Key, 0, len(l.Keys)),
	}

	for _, k := range l.Keys {
		sc, err := encodeSC(k.SC)
		if err != nil {
			return nil, fmt.Errorf("kb: layout %q: %w", l.ID, err)
		}
		vk, err := encodeVK(k.VK)
		if err != nil {
			return nil, fmt.Errorf("kb: layout %q key %s: %w", l.ID, k.SC, err)
		}

		key := &kbpb.Key{
			Sc:      sc,
			Vk:      vk,
			Name:    k.Name,
			Outputs: make([]*kbpb.Output, 0, len(k.out)),
		}
		// Mods orders the states, so the encoding is deterministic and a
		// rebuild that changed nothing produces an identical file.
		for _, m := range k.Mods() {
			o := k.out[m]
			key.Outputs = append(key.Outputs, &kbpb.Output{
				Mod:  uint32(m),
				Text: o.Text,
				Dead: o.Dead,
			})
		}

		// Mods only reports the states this package knows how to name, so a
		// key holding anything else would be written out short. Say so
		// rather than quietly losing it.
		if len(key.Outputs) != len(k.out) {
			for m := range k.out {
				if _, named := modNames[m]; !named {
					return nil, fmt.Errorf("kb: layout %q key %s: unknown modifier state %d",
						l.ID, k.SC, uint8(m))
				}
			}
		}

		msg.Keys = append(msg.Keys, key)
	}

	return proto.MarshalOptions{Deterministic: true}.Marshal(msg)
}

// Unmarshal decodes a layout from the dataset's storage format and indexes it,
// which positions the keys and discards any that sit outside the alphanumeric
// block.
func (l *Layout) Unmarshal(raw []byte) error {
	var msg kbpb.Layout
	if err := proto.Unmarshal(raw, &msg); err != nil {
		return err
	}

	form, ok := formFromProto[msg.Form]
	if !ok {
		return fmt.Errorf("kb: layout %q has unknown form %v", msg.Id, msg.Form)
	}

	*l = Layout{
		ID:      msg.Id,
		Name:    msg.Name,
		File:    msg.File,
		Form:    form,
		Locales: decodeLocales(msg.Locales),
		Keys:    make([]Key, 0, len(msg.Keys)),
	}

	for _, k := range msg.Keys {
		sc := decodeSC(k.Sc)

		vk, err := decodeVK(k.Vk)
		if err != nil {
			return fmt.Errorf("kb: layout %q key %s: %w", msg.Id, sc, err)
		}

		key := NewKey(sc, vk, k.Name)
		for _, o := range k.Outputs {
			mod := Mod(o.Mod)
			if _, named := modNames[mod]; !named {
				return fmt.Errorf("kb: layout %q key %s: unknown modifier state %d", msg.Id, sc, o.Mod)
			}
			key.Set(mod, Out{Text: o.Text, Dead: o.Dead})
		}
		l.Keys = append(l.Keys, key)
	}

	l.index()

	return nil
}

// MarshalCatalogue encodes the index of available layouts. The generator uses
// it; reading the dataset goes through List and Get.
func MarshalCatalogue(entries []Entry) ([]byte, error) {
	msg := &kbpb.Catalogue{Entries: make([]*kbpb.Entry, 0, len(entries))}

	sorted := make([]Entry, len(entries))
	copy(sorted, entries)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].ID < sorted[j].ID })

	for _, e := range sorted {
		form, ok := formToProto[e.Form]
		if !ok {
			return nil, fmt.Errorf("kb: layout %q has unknown form %q", e.ID, e.Form)
		}

		locales, err := encodeLocales(e.Locales)
		if err != nil {
			return nil, fmt.Errorf("kb: layout %q: %w", e.ID, err)
		}

		msg.Entries = append(msg.Entries, &kbpb.Entry{
			Id:      e.ID,
			Name:    e.Name,
			Form:    form,
			Locales: locales,
		})
	}

	return proto.MarshalOptions{Deterministic: true}.Marshal(msg)
}

// UnmarshalCatalogue decodes the index of available layouts.
func UnmarshalCatalogue(raw []byte) ([]Entry, error) {
	var msg kbpb.Catalogue
	if err := proto.Unmarshal(raw, &msg); err != nil {
		return nil, err
	}

	entries := make([]Entry, 0, len(msg.Entries))
	for _, e := range msg.Entries {
		form, ok := formFromProto[e.Form]
		if !ok {
			return nil, fmt.Errorf("kb: layout %q has unknown form %v", e.Id, e.Form)
		}

		entries = append(entries, Entry{
			ID:      e.Id,
			Name:    e.Name,
			Form:    form,
			Locales: decodeLocales(e.Locales),
		})
	}

	return entries, nil
}
