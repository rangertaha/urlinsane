// Copyright 2024 Rangertaha. All Rights Reserved.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.
package kb_test

import (
	"fmt"

	"github.com/rangertaha/urlinsane/pkg/kb"
)

func ExampleGet() {
	// A layout answers to its driver name, to the file that driver ships
	// in, and to any of the KLIDs Windows installs it under.
	for _, id := range []string{"kbdus", "KBDUS.DLL", "00000409", "409", "00000804"} {
		l := kb.MustGet(id)
		fmt.Printf("%-12s %s\n", id, l.ID)
	}
	// Output:
	// kbdus        kbdus
	// KBDUS.DLL    kbdus
	// 00000409     kbdus
	// 409          kbdus
	// 00000804     kbdus
}

func ExampleLayout_Adjacent() {
	us := kb.MustGet("kbdus")

	// Nearest first. "d" comes before "s" because the rows of a keyboard
	// are staggered: "e" overhangs "d" by a quarter key and "s" by three
	// quarters.
	fmt.Println(us.Adjacent("e"))

	// Shift travels with the character.
	fmt.Println(us.Adjacent("E"))

	// Output:
	// [w r d 3 4 s]
	// [W R D # $ S]
}

func ExampleLayout_Adjacent_followsTheLayout() {
	// The answer is read off the layout rather than assumed from QWERTY.
	fmt.Println(kb.MustGet("kbdgr").Adjacent("z")) // QWERTZ
	fmt.Println(kb.MustGet("kbdfr").Adjacent("a")) // AZERTY
	fmt.Println(kb.MustGet("kbddv").Adjacent("e")) // Dvorak
	// Output:
	// [t u h 6 7 g]
	// [z q & é]
	// [o u . q j p]
}

func ExampleLayout_With() {
	// An ISO board carries a key below the left Shift that an ANSI board
	// does not, and it neighbours "a".
	uk := kb.MustGet("kbduk")

	fmt.Println(uk.Adjacent("a"))
	fmt.Println(uk.With(kb.ANSI).Adjacent("a"))
	// Output:
	// [s q \ z w]
	// [s q z w]
}

func ExampleLayout_Shifted() {
	us := kb.MustGet("kbdus")

	fmt.Println(us.Shifted("4"))
	fmt.Println(us.Unshifted("$"))

	// AltGr is where the extra characters of a European layout live.
	fmt.Println(kb.MustGet("kbdfr").AltGraphed("à"))
	// Output:
	// [$]
	// [4]
	// [@]
}

func ExampleLayout_Translate() {
	us, ru, fr := kb.MustGet("kbdus"), kb.MustGet("kbdru"), kb.MustGet("kbdfr")

	// What someone types with the wrong layout selected.
	fmt.Println(us.Translate("hello", ru))
	fmt.Println(us.Translate("google.com", ru))

	// AZERTY moves the punctuation, which is what makes this a squatting
	// vector rather than a curiosity.
	fmt.Println(us.Translate("google.com", fr))
	// Output:
	// руддщ
	// пщщпдуюсщь
	// google:co,
}

func ExampleLayout_Strokes() {
	us, ru := kb.MustGet("kbdus"), kb.MustGet("kbdru")

	// Strokes are keys and modifiers, with no text of their own. Reading
	// them on another layout is what Translate does in one step.
	fmt.Println(us.Strokes("abc"))
	fmt.Println(ru.Type(us.Strokes("abc")))
	// Output:
	// [{1E base} {30 base} {2E base}]
	// фис
}

func ExampleByLanguage() {
	// A bare subtag matches every layout for the language; a full tag
	// matches only those carrying that exact locale.
	for _, e := range kb.ByLanguage("de") {
		fmt.Println(e.ID, e.Name)
	}
	fmt.Println("---")
	for _, e := range kb.ByLanguage("de-CH") {
		fmt.Println(e.ID, e.Name)
	}
	// Output:
	// kbdgr German
	// kbdgr1 German (IBM)
	// kbdgre1 German Extended (E1)
	// kbdgre2 German Extended (E2)
	// kbdsg Swiss German
	// ---
	// kbdsg Swiss German
}

func ExampleByString() {
	// Which keyboards have a key for every character of this string?
	fmt.Println(len(kb.ByString("google.com")))
	fmt.Println(len(kb.ByString("münchen")))
	fmt.Println(len(kb.ByString("中文")))
	// Output:
	// 98
	// 20
	// 0
}

func ExampleLayout_Print() {
	// Boxes are placed from the real geometry, so the indentation between
	// rows is the physical stagger. Dead keys are bracketed, since they
	// type nothing on their own.
	kb.MustGet("kbdgr").Print()
	// Output:
	// kbdgr — German (iso)
	// +---+---+---+---+---+---+---+---+---+---+---+---+---+
	// | ° | ! | " | § | $ | % | & | / | ( | ) | = | ? |[`]|
	// |[^]| 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | 0 | ß |[´]|
	// +-----+---+---+---+---+---+---+---+---+---+---+---+---+
	// |     | Q | W | E | R | T | Z | U | I | O | P | Ü | * |
	// |     | q | w | e | r | t | z | u | i | o | p | ü | + |
	// +------+---+---+---+---+---+---+---+---+---+---+---+---+
	// |      | A | S | D | F | G | H | J | K | L | Ö | Ä | ' |
	// |      | a | s | d | f | g | h | j | k | l | ö | ä | # |
	// +----+---+---+---+---+---+---+---+---+---+---+---+-+---+
	// |    | > | Y | X | C | V | B | N | M | ; | : | _ |
	// |    | < | y | x | c | v | b | n | m | , | . | - |
	// +----+---+---+-+------------------------++---+---+
	//                |                        |
	//                |                        |
	//                +------------------------+
}

func ExampleKey() {
	// A key knows the switch it sits on and what it types in each state.
	k, _ := kb.MustGet("kbdgr").Key("29")

	fmt.Println(k.SC, k.VK)

	o, _ := k.Text(kb.Base)
	fmt.Printf("%q dead=%v\n", o.Text, o.Dead)
	fmt.Printf("%q\n", k.Shift())
	// Output:
	// 29 VK_OEM_5
	// "^" dead=true
	// "°"
}

func ExampleLanguages() {
	langs := kb.Languages()
	fmt.Println(len(langs), "languages")
	fmt.Println(len(kb.List()), "layouts")
	// Output:
	// 110 languages
	// 203 layouts
}
