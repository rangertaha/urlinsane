package aci

// // adjacentCharacterInsertionFunc are created by inserting letters adjacent of each letter. For example, www.googhle.com
// // and www.goopgle.com
// func adjacentCharacterInsertionFunc(tc Result) (results []Result) {
// 	for _, keyboard := range tc.Keyboards {
// 		for i, char := range tc.Original.Domain {
// 			for _, key := range keyboard.Adjacent(string(char)) {
// 				d1 := tc.Original.Domain[:i] + string(key) + string(char) + tc.Original.Domain[i+1:]
// 				dm1 := Domain{tc.Original.Subdomain, d1, tc.Original.Suffix, Meta{}, false}
// 				results = append(results, Result{Original: tc.Original, Variant: dm1, Typo: tc.Typo, Data: tc.Data})

// 				d2 := tc.Original.Domain[:i] + string(char) + string(key) + tc.Original.Domain[i+1:]
// 				dm2 := Domain{tc.Original.Subdomain, d2, tc.Original.Suffix, Meta{}, false}
// 				results = append(results, Result{Original: tc.Original, Variant: dm2, Typo: tc.Typo, Data: tc.Data})
// 			}
// 		}
// 	}
// 	return
// }


import (
	"github.com/rangertaha/urlinsane"
	"github.com/rangertaha/urlinsane/plugins/algorithms"
)

const CODE = "aci"

type AdjacentCharacterInsertion struct {
	types []string
}

func (n *AdjacentCharacterInsertion) Id() string {
	return CODE
}

func (n *AdjacentCharacterInsertion) IsType(str string) bool {
	return algorithms.IsType(n.types, str)
}

func (n *AdjacentCharacterInsertion) Name() string {
	return "Adjacent Character Insertion"
}

func (n *AdjacentCharacterInsertion) Description() string {
	return "Adjacent Character Insertion inserts adjacent character"
}

func (n *AdjacentCharacterInsertion) Fields() []string {
	return []string{}
}

func (n *AdjacentCharacterInsertion) Headers() []string {
	return []string{}
}

func (n *AdjacentCharacterInsertion) Exec(in urlinsane.Typo) (out []urlinsane.Typo) {
	out = append(out, in)
	return
}

// Register the plugin
func init() {
	algorithms.Add(CODE, func() urlinsane.Algorithm {
		return &AdjacentCharacterInsertion{
			types: []string{algorithms.ENTITY, algorithms.DOMAIN},
		}
	})
}
