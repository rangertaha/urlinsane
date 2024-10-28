package dcr

// doubleCharacterReplacementFunc are created by replacing identical, consecutive
// letters of the domain name with adjacent letters on the keyboard.
// For example, www.gppgle.com and www.giigle.com
// func doubleCharacterReplacementFunc(tc Result) (results []Result) {
// 	for _, keyboard := range tc.Keyboards {
// 		for i, char := range tc.Original.Domain {
// 			if i < len(tc.Original.Domain)-1 {
// 				if tc.Original.Domain[i] == tc.Original.Domain[i+1] {
// 					for _, key := range keyboard.Adjacent(string(char)) {
// 						domain := tc.Original.Domain[:i] + string(key) + string(key) + tc.Original.Domain[i+2:]
// 						dm := Domain{tc.Original.Subdomain, domain, tc.Original.Suffix, Meta{}, false}
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

const CODE = "dcr"

type DoubleCharacterReplacement struct {
	types []string
}

func (n *DoubleCharacterReplacement) Id() string {
	return CODE
}
func (n *DoubleCharacterReplacement) IsType(str string) bool {
	return algorithms.IsType(n.types, str)
}

func (n *DoubleCharacterReplacement) Name() string {
	return "DoubleCharacterReplacement"
}

func (n *DoubleCharacterReplacement) Description() string {
	return "Double Character Replacement repeats a character twice"
}

func (n *DoubleCharacterReplacement) Fields() []string {
	return []string{}
}

func (n *DoubleCharacterReplacement) Headers() []string {
	return []string{}
}

func (n *DoubleCharacterReplacement) Exec(in urlinsane.Typo) (out []urlinsane.Typo) {
	out = append(out, in)
	return
}

// Register the plugin
func init() {
	algorithms.Add(CODE, func() urlinsane.Algorithm {
		return &DoubleCharacterReplacement{
			types: []string{algorithms.ENTITY, algorithms.DOMAIN},
		}
	})
}
