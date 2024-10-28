package ns

// // numeralSwapFunc are created by swapping numbers and corresponding words
// func numeralSwapFunc(tc Result) (results []Result) {
// 	for _, keyboard := range tc.Keyboards {
// 		for inum, words := range keyboard.Language.Numerals {
// 			for _, snum := range words {
// 				{
// 					dnew := strings.Replace(tc.Original.Domain, snum, inum, -1)
// 					dm := Domain{tc.Original.Subdomain, dnew, tc.Original.Suffix, Meta{}, false}
// 					if dnew != tc.Original.Domain {
// 						results = append(results, Result{Original: tc.Original, Variant: dm, Typo: tc.Typo, Data: tc.Data})
// 					}
// 				}
// 				{
// 					dnew := strings.Replace(tc.Original.Domain, inum, snum, -1)
// 					dm := Domain{tc.Original.Subdomain, dnew, tc.Original.Suffix, Meta{}, false}
// 					if dnew != tc.Original.Domain {
// 						results = append(results, Result{Original: tc.Original, Variant: dm, Typo: tc.Typo, Data: tc.Data})
// 					}
// 				}
// 			}
// 		}
// 	}
// 	return
// }


import (
	"github.com/rangertaha/urlinsane"
	"github.com/rangertaha/urlinsane/plugins/algorithms"
)

const CODE = "ns"

type NumeralSwap struct {
	types []string
}

func (n *NumeralSwap) Id() string {
	return CODE
}
func (n *NumeralSwap) IsType(str string) bool {
	return algorithms.IsType(n.types, str)
}

func (n *NumeralSwap) Name() string {
	return "NumeralSwap"
}

func (n *NumeralSwap) Description() string {
	return "Numeral Swap numbers, words and vice versa"
}

func (n *NumeralSwap) Fields() []string {
	return []string{}
}

func (n *NumeralSwap) Headers() []string {
	return []string{}
}

func (n *NumeralSwap) Exec(in urlinsane.Typo) (out []urlinsane.Typo) {
	out = append(out, in)
	return
}

// Register the plugin
func init() {
	algorithms.Add(CODE, func() urlinsane.Algorithm {
		return &NumeralSwap{
			types: []string{algorithms.ENTITY, algorithms.DOMAIN},
		}
	})
}
