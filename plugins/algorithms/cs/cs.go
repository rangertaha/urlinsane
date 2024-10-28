package cs

// // characterSwapFunc typos are when two consecutive characters are swapped in the original domain name.
// // Example: www.examlpe.com
// func characterSwapFunc(tc Result) (results []Result) {
// 	for i := range tc.Original.Domain {
// 		if i <= len(tc.Original.Domain)-2 {
// 			domain := fmt.Sprint(
// 				tc.Original.Domain[:i],
// 				string(tc.Original.Domain[i+1]),
// 				string(tc.Original.Domain[i]),
// 				tc.Original.Domain[i+2:],
// 			)
// 			if tc.Original.Domain != domain {
// 				dm := Domain{tc.Original.Subdomain, domain, tc.Original.Suffix, Meta{}, false}
// 				results = append(results, Result{Original: tc.Original, Variant: dm, Typo: tc.Typo, Data: tc.Data})
// 			}
// 		}
// 	}
// 	return results
// }

import (
	"github.com/rangertaha/urlinsane"
	"github.com/rangertaha/urlinsane/plugins/algorithms"
)

const CODE = "cs"

type CharacterSwap struct {
	types []string
}

func (n *CharacterSwap) Id() string {
	return CODE
}
func (n *CharacterSwap) IsType(str string) bool {
	return algorithms.IsType(n.types, str)
}

func (n *CharacterSwap) Name() string {
	return "Character Swap"
}

func (n *CharacterSwap) Description() string {
	return "Character Swap Swapping two consecutive characters in a domain"
}

func (n *CharacterSwap) Fields() []string {
	return []string{}
}

func (n *CharacterSwap) Headers() []string {
	return []string{}
}

func (n *CharacterSwap) Exec(in urlinsane.Typo) (out []urlinsane.Typo) {
	out = append(out, in)
	return
}

// Register the plugin
func init() {
	algorithms.Add(CODE, func() urlinsane.Algorithm {
		return &CharacterSwap{
			types: []string{algorithms.ENTITY, algorithms.DOMAIN},
		}
	})
}
