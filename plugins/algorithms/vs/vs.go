package vs

// vowelSwappingFunc swaps vowels within the domain name except for the first letter.
// For example, www.google.com becomes www.gaagle.com.
// func vowelSwappingFunc(tc Result) (results []Result) {
// 	for _, keyboard := range tc.Keyboards {
// 		for _, vchar := range keyboard.Language.Vowels {
// 			if strings.Contains(tc.Original.Domain, vchar) {
// 				for _, vvchar := range keyboard.Language.Vowels {
// 					new := strings.Replace(tc.Original.Domain, vchar, vvchar, -1)
// 					if new != tc.Original.Domain {
// 						dm := Domain{tc.Original.Subdomain, new, tc.Original.Suffix, Meta{}, false}
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

const CODE = "vs"

type VowelSwapping struct {
	types []string
}

func (n *VowelSwapping) Id() string {
	return CODE
}
func (n *VowelSwapping) IsType(str string) bool {
	return algorithms.IsType(n.types, str)
}

func (n *VowelSwapping) Name() string {
	return "Vowel Swapping"
}

func (n *VowelSwapping) Description() string {
	return "Vowel Swapping is created by swaps vowels"
}

func (n *VowelSwapping) Fields() []string {
	return []string{}
}

func (n *VowelSwapping) Headers() []string {
	return []string{}
}

func (n *VowelSwapping) Exec(in urlinsane.Typo) (out []urlinsane.Typo) {
	out = append(out, in)
	return
}

// Register the plugin
func init() {
	algorithms.Add(CODE, func() urlinsane.Algorithm {
		return &VowelSwapping{
			types: []string{algorithms.ENTITY, algorithms.DOMAIN},
		}
	})
}
