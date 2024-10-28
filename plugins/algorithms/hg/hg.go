package hg

// homoglyphFunc when one or more characters that look similar to another
// character but are different are called homogylphs. An example is that the
// lower case l looks similar to the numeral one, e.g. l vs 1. For example,
// google.com becomes goog1e.com.
// func homoglyphFunc(tc Result) (results []Result) {
// 	for i, char := range tc.Original.Domain {
// 		// Check the alphabet of the language associated with the keyboard for
// 		// homoglyphs
// 		for _, keyboard := range tc.Keyboards {
// 			for _, kchar := range keyboard.Language.SimilarChars(string(char)) {
// 				domain := fmt.Sprint(tc.Original.Domain[:i], kchar, tc.Original.Domain[i+1:])
// 				if tc.Original.Domain != domain {
// 					dm := Domain{tc.Original.Subdomain, domain, tc.Original.Suffix, Meta{}, false}
// 					results = append(results, Result{Original: tc.Original, Variant: dm, Typo: tc.Typo, Data: tc.Data})
// 				}
// 			}
// 		}
// 		// Check languages given with the (-l --language) CLI options for homoglyphs.
// 		for _, language := range tc.Languages {
// 			for _, lchar := range language.SimilarChars(string(char)) {
// 				domain := fmt.Sprint(tc.Original.Domain[:i], lchar, tc.Original.Domain[i+1:])
// 				if tc.Original.Domain != domain {
// 					dm := Domain{tc.Original.Subdomain, domain, tc.Original.Suffix, Meta{}, false}
// 					results = append(results, Result{Original: tc.Original, Variant: dm, Typo: tc.Typo, Data: tc.Data})
// 				}
// 			}
// 		}

// 	}
// 	return results
// }

import (
	"github.com/rangertaha/urlinsane"
	"github.com/rangertaha/urlinsane/plugins/algorithms"
)

const CODE = "hg"

type Homoglyphs struct {
	types []string
}

func (n *Homoglyphs) Id() string {
	return CODE
}
func (n *Homoglyphs) IsType(str string) bool {
	return algorithms.IsType(n.types, str)
}

func (n *Homoglyphs) Name() string {
	return "Homoglyphs"
}

func (n *Homoglyphs) Description() string {
	return "Replaces characters with characters that look similar"
}

func (n *Homoglyphs) Fields() []string {
	return []string{}
}

func (n *Homoglyphs) Headers() []string {
	return []string{}
}

func (n *Homoglyphs) Exec(in urlinsane.Typo) (out []urlinsane.Typo) {
	out = append(out, in)
	return
}

// Register the plugin
func init() {
	algorithms.Add(CODE, func() urlinsane.Algorithm {
		return &Homoglyphs{
			types: []string{algorithms.ENTITY, algorithms.DOMAIN},
		}
	})
}
