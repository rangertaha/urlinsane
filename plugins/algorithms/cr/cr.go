package cr

// characterRepeatFunc are created by repeating a letter of the domain name.
// Example, www.ggoogle.com and www.gooogle.com
// func characterRepeatFunc(tc Result) (results []Result) {
// 	for i := range tc.Original.Domain {
// 		if i <= len(tc.Original.Domain) {
// 			domain := fmt.Sprint(
// 				tc.Original.Domain[:i],
// 				string(tc.Original.Domain[i]),
// 				string(tc.Original.Domain[i]),
// 				tc.Original.Domain[i+1:],
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

const CODE = "cr"

type CharacterRepeat struct {
	types []string
}

func (n *CharacterRepeat) Id() string {
	return CODE
}
func (n *CharacterRepeat) IsType(str string) bool {
	return algorithms.IsType(n.types, str)
}

func (n *CharacterRepeat) Name() string {
	return "CharacterRepeat"
}

func (n *CharacterRepeat) Description() string {
	return "Character Repeat Repeats a character of the domain name twice"
}

func (n *CharacterRepeat) Fields() []string {
	return []string{}
}

func (n *CharacterRepeat) Headers() []string {
	return []string{}
}

func (n *CharacterRepeat) Exec(in urlinsane.Typo) (out []urlinsane.Typo) {
	out = append(out, in)
	return
}

// Register the plugin
func init() {
	algorithms.Add(CODE, func() urlinsane.Algorithm {
		return &CharacterRepeat{
			types: []string{algorithms.ENTITY, algorithms.DOMAIN},
		}
	})
}
