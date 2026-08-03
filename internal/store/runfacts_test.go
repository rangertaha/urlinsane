// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package store_test

import (
	"reflect"
	"testing"

	"github.com/rangertaha/urlinsane/internal/plugins/report"
	"github.com/rangertaha/urlinsane/internal/store"
)

// notPersisted are the report.Options fields that legitimately do not survive a
// save, with the reason each one does not.
//
// Everything else is a fact about the run that produced the graph, and a run
// fact that is not in the index is a fact `report` renders as its zero value —
// off the very bytes that had it.
var notPersisted = map[string]string{
	"Target":  "stored as Entry.Key, the canonical seed",
	"Filters": "a flag on the render, not a fact about the scan",
	"Verbose": "a flag on the render",
	"Color":   "a property of the terminal being written to",
	"Format":  "chosen per render by -o",
	"Elapsed": "deliberately excluded from canonical output; it differs every run",
	"Partial": "stored as Entry.Partial",

	"PartialWhy": "derived from Entry.Partial by the caller",
}

// A run fact the renderer reads must be carried by the index.
//
// This is a class test, not a test of three fields. The scan root is content
// addressed and deliberately carries no invocation detail — that is what makes
// two identical scans share a CID — so every such fact has to live in the Entry
// beside it. Nothing enforced that, and it went wrong twice: first the observer
// set, which made `typo --save-graph` say absent:8 and `report` say live:8
// (fixed in 9944703), then Scope, Plan and Rounds, which `report` rendered as
// empty and 0.
//
// Both were invisible until someone compared two renderings by eye. So the rule
// is enforced here rather than remembered: add a field to report.Options that
// describes the run, and this fails until the index carries it too.
func TestEveryRunFactSurvivesTheIndex(t *testing.T) {
	opts := reflect.TypeOf(report.Options{})
	entry := reflect.TypeOf(store.Entry{})

	carried := make(map[string]bool, entry.NumField())
	for i := 0; i < entry.NumField(); i++ {
		carried[entry.Field(i).Name] = true
	}

	for i := 0; i < opts.NumField(); i++ {
		name := opts.Field(i).Name
		if _, ok := notPersisted[name]; ok {
			continue
		}
		if !carried[name] {
			t.Errorf("report.Options.%s describes the run but store.Entry does not carry it, "+
				"so `report` renders it as its zero value off bytes that had it. "+
				"Either add it to Entry (and to the save and report paths) or, if it is a "+
				"property of the render rather than of the scan, list it in notPersisted with "+
				"the reason.", name)
		}
	}
}

// The allowlist must not rot into a way of silencing the test: every name in it
// has to still be a field.
func TestNotPersistedNamesRealFields(t *testing.T) {
	opts := reflect.TypeOf(report.Options{})
	for name := range notPersisted {
		if _, ok := opts.FieldByName(name); !ok {
			t.Errorf("notPersisted lists %q, which is no longer a report.Options field", name)
		}
	}
}
