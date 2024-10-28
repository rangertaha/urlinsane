package acs

// adjacentCharacterSubstitutionFunc typos are when characters are replaced in the original domain name by their
// adjacent ones on a specific keyboard layout, e.g., www.ezample.com, where “x” was replaced by the adjacent
// character “z” in a the QWERTY keyboard layout.
// func adjacentCharacterSubstitutionFunc(tc Result) (results []Result) {
// 	for _, keyboard := range tc.Keyboards {
// 		for i, char := range tc.Original.Domain {
// 			for _, key := range keyboard.Adjacent(string(char)) {
// 				domain := tc.Original.Domain[:i] + string(key) + tc.Original.Domain[i+1:]
// 				dm := Domain{tc.Original.Subdomain, domain, tc.Original.Suffix, Meta{}, false}
// 				results = append(results, Result{Original: tc.Original, Variant: dm, Typo: tc.Typo, Data: tc.Data})
// 			}
// 		}
// 	}
// 	return
// }


import (
	"github.com/rangertaha/urlinsane"
	"github.com/rangertaha/urlinsane/plugins/algorithms"
)

const CODE = "acs"

type AdjacentCharacterSubstitution struct {
	types []string
}

func (n *AdjacentCharacterSubstitution) Id() string {
	return CODE
}

func (n *AdjacentCharacterSubstitution) IsType(str string) bool {
	return algorithms.IsType(n.types, str)
}

func (n *AdjacentCharacterSubstitution) Name() string {
	return "Adjacent Character Substitution"
}

func (n *AdjacentCharacterSubstitution) Description() string {
	return ""
}

func (n *AdjacentCharacterSubstitution) Fields() []string {
	return []string{}
}

func (n *AdjacentCharacterSubstitution) Headers() []string {
	return []string{}
}

func (n *AdjacentCharacterSubstitution) Exec(in urlinsane.Typo) (out []urlinsane.Typo) {
	out = append(out, in)
	return
}

// Register the plugin
func init() {
	algorithms.Add(CODE, func() urlinsane.Algorithm {
		return &AdjacentCharacterSubstitution{
			types: []string{algorithms.ENTITY, algorithms.DOMAIN},
		}
	})
}
