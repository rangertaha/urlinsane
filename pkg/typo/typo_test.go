// Copyright (C) 2024 Rangertaha
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
package typo

import (
	"reflect"
	"sort"
	"testing"
)

type TypoCase struct {
	name   string
	layout string
	typos  []string
}

func TestPrefixInsertion(t *testing.T) {
	tests := []TypoCase{
		{
			name: "example",
			typos: []string{
				"aexample",
				"bexample",
				"cexample",
			},
		},
		{
			name: "rangetaha@gmail.com",
			typos: []string{
				"arangetaha@gmail.com",
				"brangetaha@gmail.com",
				"crangetaha@gmail.com",
			},
		},
		{
			name: "example.com",
			typos: []string{
				"aexample.com",
				"bexample.com",
				"cexample.com",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			variants := PrefixInsertion(test.name, "a", "b", "c")
			sort.Strings(variants)

			if !reflect.DeepEqual(variants, test.typos) {
				t.Errorf("PrefixInsertion(%s, a, b, c) = %s; want %s", test.name, variants, test.typos)
			}
		})
	}

}

func TestSuffixInsertion(t *testing.T) {
	tests := []TypoCase{
		{
			name: "example",
			typos: []string{
				"examplea",
				"exampleb",
				"examplec",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			variants := SuffixInsertion(test.name, "a", "b", "c")
			sort.Strings(variants)

			if !reflect.DeepEqual(variants, test.typos) {
				t.Errorf("SuffixInsertion(%s, a, b, c) = %s; want %s", test.name, variants, test.typos)
			}
		})
	}

}

func TestCharacterSwap(t *testing.T) {
	tests := []TypoCase{
		{
			name:   "example",
			layout: "",
			typos: []string{
				"eaxmple", "examlpe", "exampel",
				"exapmle", "exmaple", "xeample",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			variants := CharacterSwap(test.name)
			sort.Strings(variants)

			if !reflect.DeepEqual(variants, test.typos) {
				t.Errorf("CharacterSwap(%s, a, b, c) = %s; want %s", test.name, variants, test.typos)
			}
		})
	}
}

func TestAdjacentCharacterSubstitution(t *testing.T) {
	tests := []TypoCase{
		{
			name:   "example",
			layout: "QWERTY",
			typos: []string{
				"3xample", "dxample", "ecample", "esample", "exajple", "exam0le",
				"examole", "exampke", "exampl3", "exampld", "examplr", "examplw",
				"exampoe", "exanple", "exqmple", "exsmple", "exzmple", "ezample",
				"rxample", "wxample",
			},
		},
		{
			name:   "example",
			layout: "AZERTY",
			typos: []string{
				"3xample", "dxample", "ecample", "esample", "ewample", "ex1mple",
				"exalple", "exam0le", "exammle", "examole", "exampke", "exampl3",
				"exampld", "examplr", "examplz", "exampme", "exampoe", "exapple",
				"exqmple", "exzmple", "rxample", "zxample",
			},
		},

		{
			name:   "example",
			layout: "QWERTZ",
			typos: []string{
				"3xample", "dxample", "ecample", "esample", "exajple", "exam0le",
				"examole", "exampke", "exampl3", "exampld", "examplr", "examplw",
				"exampoe", "exanple", "exqmple", "exsmple", "exymple", "eyample",
				"rxample", "wxample",
			},
		},
		{
			name:   "example",
			layout: "DVORAK",
			typos: []string{
				"ebample", "eiample", "ekample", "exabple", "exahple", "exam4le",
				"examp0e", "examplj", "examplo", "examplu", "exampre", "exampse",
				"examule", "examyle", "exawple", "exomple", "jxample", "oxample",
				"uxample",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for _, keyboard := range Keyboards {
				if keyboard.name == test.layout {
					variants := AdjacentCharacterSubstitution(test.name, keyboard.layout...)
					sort.Strings(variants)

					if !reflect.DeepEqual(variants, test.typos) {
						t.Errorf("AdjacentCharacterSubstitution(%s) = %s; want %s", test.name, variants, test.typos)
					}
				}

			}
		})
	}
}

func TestAdjacentCharacterInsertion(t *testing.T) {
	tests := []TypoCase{
		{
			name:   "example",
			layout: "QWERTY",
			typos: []string{
				"3example", "dexample", "e3xample", "ecxample", "edxample",
				"erxample", "esxample", "ewxample", "exajmple", "exam0ple",
				"examjple", "examnple", "examople", "examp0le", "exampkle",
				"exampl3e", "examplde", "example3", "exampled", "exampler",
				"examplew", "examplke", "examploe", "examplre", "examplwe",
				"exampole", "exampole", "exanmple", "exaqmple", "exasmple",
				"exazmple", "excample", "exqample", "exsample", "exsample",
				"exzample", "exzample", "ezxample", "rexample", "wexample",
			},
		},
		{
			name:   "example",
			layout: "AZERTY",
			typos: []string{
				"3example", "dexample", "e3xample", "ecxample", "edxample",
				"erxample", "esxample", "ewxample", "ex1ample", "exa1mple",
				"exalmple", "exam0ple", "examlple", "exammple", "examople",
				"examp0le", "exampkle", "exampl3e", "examplde", "example3",
				"exampled", "exampler", "examplez", "examplke", "examplme",
				"examploe", "examplre", "examplze", "exampmle", "exampmle",
				"exampole", "exampole", "exampple", "exapmple", "exaqmple",
				"exazmple", "excample", "exqample", "exsample", "exwample",
				"exzample", "ezxample", "rexample", "zexample",
			},
		},

		{
			name:   "example",
			layout: "QWERTZ",
			typos: []string{
				"3example", "dexample", "e3xample", "ecxample", "edxample",
				"erxample", "esxample", "ewxample", "exajmple", "exam0ple",
				"examjple", "examnple", "examople", "examp0le", "exampkle",
				"exampl3e", "examplde", "example3", "exampled", "exampler",
				"examplew", "examplke", "examploe", "examplre", "examplwe",
				"exampole", "exampole", "exanmple", "exaqmple", "exasmple",
				"exaymple", "excample", "exqample", "exsample", "exsample",
				"exyample", "exyample", "eyxample", "rexample", "wexample",
			},
		},
		{
			name:   "example",
			layout: "DVORAK",
			typos: []string{
				"ebxample", "eixample", "ejxample", "ekxample", "eoxample",
				"euxample", "exabmple", "exahmple", "exam4ple", "exambple",
				"examhple", "examp0le", "examp4le", "exampl0e", "examplej",
				"exampleo", "exampleu", "examplje", "examploe", "examplre",
				"examplse", "examplue", "examprle", "exampsle", "exampule",
				"exampyle", "examuple", "examwple", "examyple", "exaomple",
				"exawmple", "exbample", "exiample", "exkample", "exoample",
				"jexample", "oexample", "uexample",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for _, keyboard := range Keyboards {
				if keyboard.name == test.layout {
					variants := AdjacentCharacterInsertion(test.name, keyboard.layout...)
					sort.Strings(variants)

					if !reflect.DeepEqual(variants, test.typos) {
						t.Errorf("AdjacentCharacterInsertion(%s) = %s; want %s", test.name, variants, test.typos)
					}
				}

			}
		})
	}
}

func TestHyphenInsertion(t *testing.T) {
	tests := []TypoCase{
		{
			name:   "example",
			layout: "QWERTY",
			typos: []string{
				"-example", "e-xample", "ex-ample", "exa-mple",
				"exam-ple", "examp-le", "example-",
			},
		},
		{
			name:   "rangertaha",
			layout: "AZERTY",
			typos: []string{
				"-rangertaha", "r-angertaha", "ra-ngertaha", "ran-gertaha",
				"rang-ertaha", "range-rtaha", "ranger-taha", "rangert-aha",
				"rangerta-ha", "rangertaha-",
			},
		},

		{
			name:   "alessandro",
			layout: "QWERTZ",
			typos: []string{
				"-alessandro", "a-lessandro", "al-essandro", "ale-ssandro",
				"ales-sandro", "aless-andro", "alessa-ndro", "alessan-dro",
				"alessand-ro", "alessandro-",
			},
		},
		{
			name:   "puppet",
			layout: "DVORAK",
			typos: []string{
				"-puppet", "p-uppet", "pu-ppet", "pup-pet", "pupp-et", "puppet-",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for _, keyboard := range Keyboards {
				if keyboard.name == test.layout {
					variants := HyphenInsertion(test.name)
					sort.Strings(variants)

					if !reflect.DeepEqual(variants, test.typos) {
						t.Errorf("HyphenInsertion(%s) = %s; want %s", test.name, variants, test.typos)
					}
				}

			}
		})
	}
}

func TestGraphemeInsertion(t *testing.T) {
	tests := []TypoCase{
		{
			name:   "example",
			layout: "QWERTY",
			typos: []string{
				"aexample", "bexample", "cexample", "dexample", "eaxample",
				"ebxample", "ecxample", "edxample", "eexample", "eexample",
				"efxample", "egxample", "ehxample", "eixample", "ejxample",
				"ekxample", "elxample", "emxample", "enxample", "eoxample",
				"epxample", "eqxample", "erxample", "esxample", "etxample",
				"euxample", "evxample", "ewxample", "exaample", "exaample",
				"exabmple", "exacmple", "exadmple", "exaemple", "exafmple",
				"exagmple", "exahmple", "exaimple", "exajmple", "exakmple",
				"exalmple", "examaple", "exambple", "examcple", "examdple",
				"exameple", "examfple", "examgple", "examhple", "examiple",
				"examjple", "examkple", "examlple", "exammple", "exammple",
				"examnple", "examople", "exampale", "exampble", "exampcle",
				"exampdle", "exampele", "exampfle", "exampgle", "examphle",
				"exampile", "exampjle", "exampkle", "examplea", "exampleb",
				"examplec", "exampled", "examplee", "examplef", "exampleg",
				"exampleh", "examplei", "examplej", "examplek", "examplel",
				"examplem", "examplen", "exampleo", "examplep", "exampleq",
				"exampler", "examples", "examplet", "exampleu", "examplev",
				"examplew", "examplex", "exampley", "examplez", "examplle",
				"exampmle", "exampnle", "exampole", "exampple", "exampple",
				"exampqle", "examprle", "exampsle", "examptle", "exampule",
				"exampvle", "exampwle", "exampxle", "exampyle", "exampzle",
				"examqple", "examrple", "examsple", "examtple", "examuple",
				"examvple", "examwple", "examxple", "examyple", "examzple",
				"exanmple", "exaomple", "exapmple", "exaqmple", "exarmple",
				"exasmple", "exatmple", "exaumple", "exavmple", "exawmple",
				"exaxmple", "exaymple", "exazmple", "exbample", "excample",
				"exdample", "exeample", "exfample", "exgample", "exhample",
				"exiample", "exjample", "exkample", "exlample", "exmample",
				"exnample", "exoample", "expample", "exqample", "exrample",
				"exsample", "extample", "exuample", "exvample", "exwample",
				"exxample", "exxample", "exyample", "exzample", "eyxample",
				"ezxample", "fexample", "gexample", "hexample", "iexample",
				"jexample", "kexample", "lexample", "mexample", "nexample",
				"oexample", "pexample", "qexample", "rexample", "sexample",
				"texample", "uexample", "vexample", "wexample", "xexample",
				"yexample", "zexample",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for _, keyboard := range Keyboards {
				if keyboard.name == test.layout {
					variants := GraphemeInsertion(test.name, Graphemes...)
					sort.Strings(variants)

					if !reflect.DeepEqual(variants, test.typos) {
						t.Errorf("GraphemeInsertion(%s) = %s; want %s", test.name, variants, test.typos)
					}
				}

			}
		})
	}
}

func TestGraphemeReplacement(t *testing.T) {
	tests := []TypoCase{
		{
			name:   "example",
			layout: "QWERTY",
			typos: []string{
				"axample", "bxample", "cxample", "dxample", "eaample", "ebample",
				"ecample", "edample", "eeample", "efample", "egample", "ehample",
				"eiample", "ejample", "ekample", "elample", "emample", "enample",
				"eoample", "epample", "eqample", "erample", "esample", "etample",
				"euample", "evample", "ewample", "exaaple", "exabple", "exacple",
				"exadple", "exaeple", "exafple", "exagple", "exahple", "exaiple",
				"exajple", "exakple", "exalple", "examale", "examble", "examcle",
				"examdle", "examele", "examfle", "examgle", "examhle", "examile",
				"examjle", "examkle", "examlle", "exammle", "examnle", "examole",
				"exampae", "exampbe", "exampce", "exampde", "exampee", "exampfe",
				"exampge", "examphe", "exampie", "exampje", "exampke", "exampla",
				"examplb", "examplc", "exampld", "example", "example", "example",
				"example", "example", "example", "example", "examplf", "examplg",
				"examplh", "exampli", "examplj", "examplk", "exampll", "examplm",
				"exampln", "examplo", "examplp", "examplq", "examplr", "exampls",
				"examplt", "examplu", "examplv", "examplw", "examplx", "examply",
				"examplz", "exampme", "exampne", "exampoe", "examppe", "exampqe",
				"exampre", "exampse", "exampte", "exampue", "exampve", "exampwe",
				"exampxe", "exampye", "exampze", "examqle", "examrle", "examsle",
				"examtle", "examule", "examvle", "examwle", "examxle", "examyle",
				"examzle", "exanple", "exaople", "exapple", "exaqple", "exarple",
				"exasple", "exatple", "exauple", "exavple", "exawple", "exaxple",
				"exayple", "exazple", "exbmple", "excmple", "exdmple", "exemple",
				"exfmple", "exgmple", "exhmple", "eximple", "exjmple", "exkmple",
				"exlmple", "exmmple", "exnmple", "exomple", "expmple", "exqmple",
				"exrmple", "exsmple", "extmple", "exumple", "exvmple", "exwmple",
				"exxmple", "exymple", "exzmple", "eyample", "ezample", "fxample",
				"gxample", "hxample", "ixample", "jxample", "kxample", "lxample",
				"mxample", "nxample", "oxample", "pxample", "qxample", "rxample",
				"sxample", "txample", "uxample", "vxample", "wxample", "xxample",
				"yxample", "zxample",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for _, keyboard := range Keyboards {
				if keyboard.name == test.layout {
					variants := GraphemeReplacement(test.name, Graphemes...)
					sort.Strings(variants)

					if !reflect.DeepEqual(variants, test.typos) {
						t.Errorf("GraphemeReplacement(%s) = %s; want %s", test.name, variants, test.typos)
					}
				}

			}
		})
	}
}

func TestGraphemeRepetition(t *testing.T) {
	tests := []TypoCase{
		{
			name:   "example",
			layout: "QWERTY",
			typos: []string{
				"eexample", "exaample", "exammple", "examplee", "examplle",
				"exampple", "exxample",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for _, keyboard := range Keyboards {
				if keyboard.name == test.layout {
					variants := GraphemeRepetition(test.name)
					sort.Strings(variants)

					if !reflect.DeepEqual(variants, test.typos) {
						t.Errorf("GraphemeRepetition(%s) = %s; want %s", test.name, variants, test.typos)
					}
				}

			}
		})
	}
}

func TestDoubleGraphemeAdjacentReplacement(t *testing.T) {
	tests := []TypoCase{
		{
			name:   "google",
			layout: "QWERTY",
			typos: []string{
				"g99gle", "giigle", "gllgle", "gppgle",
			},
		},
		{
			name:   "facebook",
			layout: "QWERTY",
			typos: []string{
				"faceb99k", "facebiik", "facebllk", "facebppk",
			},
		},
		{
			name:   "zoom",
			layout: "QWERTY",
			typos: []string{
				"z99m", "ziim", "zllm", "zppm",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for _, keyboard := range Keyboards {
				if keyboard.name == test.layout {
					variants := DoubleGraphemeAdjacentReplacement(test.name, keyboard.layout...)
					sort.Strings(variants)

					if !reflect.DeepEqual(variants, test.typos) {
						t.Errorf("DoubleGraphemeAdjacentReplacement(%s) = %s; want %s", test.name, variants, test.typos)
					}
				}

			}
		})
	}
}

func TestGraphemeOmission(t *testing.T) {
	tests := []TypoCase{
		{
			name:   "google",
			layout: "QWERTY",
			typos: []string{
				"gogle", "gogle", "googe", "googl", "goole", "oogle",
			},
		},
		{
			name:   "facebook",
			layout: "QWERTY",
			typos: []string{
				"acebook", "facbook", "facebok", "facebok", "faceboo",
				"faceook", "faebook", "fcebook",
			},
		},
		{
			name:   "zoom",
			layout: "QWERTY",
			typos: []string{
				"oom", "zom", "zom", "zoo",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for _, keyboard := range Keyboards {
				if keyboard.name == test.layout {
					variants := GraphemeOmission(test.name)
					sort.Strings(variants)

					if !reflect.DeepEqual(variants, test.typos) {
						t.Errorf("GraphemeOmission(%s) = %s; want %s", test.name, variants, test.typos)
					}
				}

			}
		})
	}
}

func TestSingularPluraliseSubstitution(t *testing.T) {
	tests := []TypoCase{
		{
			name:   "google",
			layout: "QWERTY",
			typos: []string{
				"googles",
			},
		},
		{
			name:   "facebooking",
			layout: "QWERTY",
			typos: []string{
				"facebookings",
			},
		},
		{
			name:   "zooms",
			layout: "QWERTY",
			typos: []string{
				"zoom",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for _, keyboard := range Keyboards {
				if keyboard.name == test.layout {
					variants := SingularPluraliseSubstitution(test.name)
					sort.Strings(variants)

					if !reflect.DeepEqual(variants, test.typos) {
						t.Errorf("SingularPluraliseSubstitution(%s) = %s; want %s", test.name, variants, test.typos)
					}
				}

			}
		})
	}
}
