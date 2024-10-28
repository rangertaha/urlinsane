package cm

// CcommonMisspellingsFunc are created with common misspellings in the given
// language. For example, www.youtube.com becomes www.youtub.com and
// www.abseil.com becomes www.absail.com
// func commonMisspellingsFunc(tc Result) (results []Result) {
// 	for _, keyboard := range tc.Keyboards {
// 		for _, word := range keyboard.Language.SimilarSpellings(tc.Original.Domain) {
// 			dm := Domain{tc.Original.Subdomain, word, tc.Original.Suffix, Meta{}, false}
// 			results = append(results, Result{Original: tc.Original, Variant: dm, Typo: tc.Typo, Data: tc.Data})

// 		}
// 	}
// 	return
// }

import (
	"github.com/rangertaha/urlinsane"
	"github.com/rangertaha/urlinsane/plugins/algorithms"
)

const CODE = "cm"

type CommonMisspellings struct {
	types []string
}

func (n *CommonMisspellings) Id() string {
	return CODE
}
func (n *CommonMisspellings) IsType(str string) bool {
	return algorithms.IsType(n.types, str)
}

func (n *CommonMisspellings) Name() string {
	return "Common Misspellings"
}

func (n *CommonMisspellings) Description() string {
	return "Common Misspellings are created from a dictionary of commonly misspelled words"
}

func (n *CommonMisspellings) Fields() []string {
	return []string{}
}

func (n *CommonMisspellings) Headers() []string {
	return []string{}
}

func (n *CommonMisspellings) Exec(in urlinsane.Typo) (out []urlinsane.Typo) {
	out = append(out, in)
	return
}

// Register the plugin
func init() {
	algorithms.Add(CODE, func() urlinsane.Algorithm {
		return &CommonMisspellings{
			types: []string{algorithms.ENTITY, algorithms.DOMAIN},
		}
	})
}
