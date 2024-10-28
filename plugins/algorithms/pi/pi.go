package pi

// // adjacentCharacterInsertionFunc are created by inserting letters adjacent of each letter. For example, www.googhle.com
// // and www.goopgle.com
// func periodInsertionFunc(tc Result) (results []Result) {

// 	for i, char := range tc.Original.Domain {

// 		d1 := tc.Original.Domain[:i] + "." + string(char) + tc.Original.Domain[i+1:]
// 		dm1 := Domain{tc.Original.Subdomain, d1, tc.Original.Suffix, Meta{}, false}
// 		results = append(results, Result{Original: tc.Original, Variant: dm1, Typo: tc.Typo, Data: tc.Data})
// 	}

// 	return
// }

import (
	"github.com/rangertaha/urlinsane"
	"github.com/rangertaha/urlinsane/plugins/algorithms"
)

const CODE = "pi"

type PeriodInsertion struct {
	types []string
}

func (n *PeriodInsertion) Id() string {
	return CODE
}
func (n *PeriodInsertion) IsType(str string) bool {
	return algorithms.IsType(n.types, str)
}

func (n *PeriodInsertion) Name() string {
	return "PeriodInsertion"
}

func (n *PeriodInsertion) Description() string {
	return "Inserting periods in the target name"
}

func (n *PeriodInsertion) Fields() []string {
	return []string{}
}

func (n *PeriodInsertion) Headers() []string {
	return []string{}
}

func (n *PeriodInsertion) Exec(in urlinsane.Typo) (out []urlinsane.Typo) {
	out = append(out, in)
	return
}

// Register the plugin
func init() {
	algorithms.Add(CODE, func() urlinsane.Algorithm {
		return &PeriodInsertion{
			types: []string{algorithms.ENTITY, algorithms.DOMAIN},
		}
	})
}
