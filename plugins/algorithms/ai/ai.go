package ia

// func alphabetInsertionFunc(tc Result) (results []Result) {
// 	alphabet := map[string]bool{}
// 	for _, keyboard := range tc.Keyboards {
// 		for _, a := range keyboard.Language.Graphemes {
// 			alphabet[a] = true
// 		}
// 	}
// 	for i, char := range tc.Original.Domain {
// 		for alp := range alphabet {
// 			d1 := tc.Original.Domain[:i] + alp + string(char) + tc.Original.Domain[i+1:]
// 			if i == len(tc.Original.Domain)-1 {
// 				d1 = tc.Original.Domain[:i] + string(char) + alp + tc.Original.Domain[i+1:]
// 			}
// 			dm1 := Domain{tc.Original.Subdomain, d1, tc.Original.Suffix, Meta{}, false}
// 			results = append(results, Result{Original: tc.Original, Variant: dm1, Typo: tc.Typo, Data: tc.Data})
// 		}
// 	}
// 	return
// }

import (
	"github.com/rangertaha/urlinsane"
	"github.com/rangertaha/urlinsane/plugins/algorithms"
)

const CODE = "ia"

type AlphabetInsertion struct {
	types []string
}

func (n *AlphabetInsertion) Id() string {
	return CODE
}
func (n *AlphabetInsertion) IsType(str string) bool {
	return algorithms.IsType(n.types, str)
}

func (n *AlphabetInsertion) Name() string {
	return "Alphabet Insertion"
}

func (n *AlphabetInsertion) Description() string {
	return "Inserting the language specific alphabet in the target domain"
}

func (n *AlphabetInsertion) Fields() []string {
	return []string{}
}

func (n *AlphabetInsertion) Headers() []string {
	return []string{}
}

func (n *AlphabetInsertion) Exec(in urlinsane.Typo) (out []urlinsane.Typo) {
	out = append(out, in)
	return
}

// Register the plugin
func init() {
	algorithms.Add(CODE, func() urlinsane.Algorithm {
		return &AlphabetInsertion{
			types: []string{algorithms.ENTITY, algorithms.DOMAIN},
		}
	})
}
