// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later
package typo

import (
	"reflect"
	"sort"
	"testing"
)

func TestTopLevelDomain(t *testing.T) {
	tests := []TypoCase{
		{
			name: "co",
			typos: []string{
				"co",
				"io",
				"uk",
			},
		},
		{
			name: "uk.com",
			typos: []string{
				"co",
				"io",
				"uk",
			},
		},
		{
			name: "uk.eu.org",
			typos: []string{
				"co",
				"io",
				"uk",
			},
		},
		{
			name: "io",
			typos: []string{
				"co",
				"io",
				"uk",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			variants := TopLevelDomain(test.name, tstTLDs...)
			sort.Strings(variants)

			if !reflect.DeepEqual(variants, test.typos) {
				t.Errorf("TopLevelDomain(%s) = %s; want %s", test.name, variants, test.typos)
			}

		})
	}
}

func TestTopLevel2Domain(t *testing.T) {
	tests := []TypoCase{
		{
			name: "co",
			typos: []string{
				"uk.com",
				"uk.io",
			},
		},
		{
			name: "uk.com",
			typos: []string{
				"uk.com",
				"uk.io",
			},
		},
		{
			name: "uk.eu.org",
			typos: []string{
				"uk.com",
				"uk.io",
			},
		},
		{
			name: "io",
			typos: []string{
				"uk.com",
				"uk.io",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			variants := SecondLevelDomain(test.name, tstTLDs...)
			sort.Strings(variants)

			if !reflect.DeepEqual(variants, test.typos) {
				t.Errorf("SecondLevelDomain(%s) = %s; want %s", test.name, variants, test.typos)
			}

		})
	}
}

func TestTopLevel3Domain(t *testing.T) {
	tests := []TypoCase{
		{
			name: "co",
			typos: []string{
				"uk.eu.org",
			},
		},
		{
			name: "uk.com",
			typos: []string{
				"uk.eu.org",
			},
		},
		{
			name: "uk.eu.org",
			typos: []string{
				"uk.eu.org",
			},
		},
		{
			name: "io",
			typos: []string{
				"uk.eu.org",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			variants := ThirdLevelDomain(test.name, tstTLDs...)
			sort.Strings(variants)

			if !reflect.DeepEqual(variants, test.typos) {
				t.Errorf("ThirdLevelDomain(%s) = %s; want %s", test.name, variants, test.typos)
			}

		})
	}
}
