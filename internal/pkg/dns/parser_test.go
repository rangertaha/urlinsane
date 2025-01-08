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
package dns

import (
	"testing"
)

var ParseTests = []TestCase{
	{
		name: "example.com",
		parts: []string{
			"",
			"example",
			"com",
		},
	},
	{
		name: "example.schools.nsw.edu.au",
		parts: []string{
			"",
			"example",
			"schools.nsw.edu.au",
		},
	},
	{
		name: "www.example.schools.nsw.edu.au",
		parts: []string{
			"www",
			"example",
			"schools.nsw.edu.au",
		},
	},
	{
		name: "c-n7k-n04-01.rz.example.com",
		parts: []string{
			"c-n7k-n04-01.rz",
			"example",
			"com",
		},
	},
	{
		name: "www.rebecca.users.example.com",
		parts: []string{
			"www.rebecca.users",
			"example",
			"com",
		},
	},
}

func TestParse(t *testing.T) {
	// parser := New()
	// d := parser.Parse("google.com")
	// // fmt.Println(d)
	// t.Error(d)

	// for _, test := range TESTS {
	// 	t.Run(test.name, func(t *testing.T) {
	// 		d := parser.Parse(test.name)
	// 		if !reflect.DeepEqual([]string{d.Prefix, d.Name, d.Suffix}, test.parts) {
	// 			t.Errorf("Split(%s) = %s, %s, %s; want %s, %s, %s", test.name, d.Prefix, d.Name, d.Suffix, test.parts[0], test.parts[1], test.parts[2])
	// 		}
	// 	})
	// }

}

// func TestJoin(t *testing.T) {

// 	for _, test := range TESTS {
// 		t.Run(test.name, func(t *testing.T) {
// 			name := Join(test.parts...)
// 			if !reflect.DeepEqual(name, test.name) {
// 				t.Errorf("Join(%s, %s, %s) = %s; want %s", test.parts[0], test.parts[1], test.parts[2], name, test.name)
// 			}
// 		})
// 	}

// }
