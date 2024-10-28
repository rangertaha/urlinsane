package ar

// func alphabetReplacementnFunc(tc Result) (results []Result) {
// 	alphabet := map[string]bool{}
// 	for _, keyboard := range tc.Keyboards {
// 		for _, a := range keyboard.Language.Graphemes {
// 			alphabet[a] = true
// 		}
// 	}

// 	for i := range tc.Original.Domain {
// 		for alp := range alphabet {
// 			d1 := tc.Original.Domain[:i] + alp + tc.Original.Domain[i+1:]
// 			dm1 := Domain{tc.Original.Subdomain, d1, tc.Original.Suffix, Meta{}, false}
// 			results = append(results, Result{Original: tc.Original, Variant: dm1, Typo: tc.Typo, Data: tc.Data})

// 			if i == len(tc.Original.Domain)-1 {
// 				d1 = tc.Original.Domain[:i] + alp + tc.Original.Domain[i+1:]
// 				dm1 := Domain{tc.Original.Subdomain, d1, tc.Original.Suffix, Meta{}, false}
// 				results = append(results, Result{Original: tc.Original, Variant: dm1, Typo: tc.Typo, Data: tc.Data})
// 			}
// 		}
// 	}
// 	return
// }

import (
	"github.com/rangertaha/urlinsane"
	"github.com/rangertaha/urlinsane/plugins/algorithms"
)

const CODE = "ar"

type AlphabetReplacement struct {
	types []string
}

func (n *AlphabetReplacement) Id() string {
	return CODE
}
func (n *AlphabetReplacement) IsType(str string) bool {
	return algorithms.IsType(n.types, str)
}

func (n *AlphabetReplacement) Name() string {
	return "Alphabet Replacement"
}

func (n *AlphabetReplacement) Description() string {
	return "Replaces an alphabet in the target domain"
}

func (n *AlphabetReplacement) Fields() []string {
	return []string{}
}

func (n *AlphabetReplacement) Headers() []string {
	return []string{}
}

func (n *AlphabetReplacement) Exec(in urlinsane.Typo) (out []urlinsane.Typo) {
	out = append(out, in)
	return
}

// Register the plugin
func init() {
	algorithms.Add(CODE, func() urlinsane.Algorithm {
		return &AlphabetReplacement{

			types: []string{algorithms.ENTITY, algorithms.DOMAIN},
		}
	})
}
