package hi

// adjacentCharacterInsertionFunc are created by inserting letters adjacent of each letter. For example, www.googhle.com
// and www.goopgle.com
// func hyphenInsertionFunc(tc Result) (results []Result) {

// 	for i, char := range tc.Original.Domain {
// 		d1 := tc.Original.Domain[:i] + "-" + string(char) + tc.Original.Domain[i+1:]
// 		if i == len(tc.Original.Domain)-1 {
// 			d1 = tc.Original.Domain[:i] + string(char) + "-" + tc.Original.Domain[i+1:]
// 		}
// 		dm1 := Domain{tc.Original.Subdomain, d1, tc.Original.Suffix, Meta{}, false}
// 		results = append(results, Result{Original: tc.Original, Variant: dm1, Typo: tc.Typo, Data: tc.Data})
// 	}
// 	return
// }

import (
	"github.com/rangertaha/urlinsane"
	"github.com/rangertaha/urlinsane/plugins/algorithms"
)

const CODE = "hi"

type DashInsertion struct {
	types []string
}

func (n *DashInsertion) Id() string {
	return CODE
}
func (n *DashInsertion) IsType(str string) bool {
	return algorithms.IsType(n.types, str)
}

func (n *DashInsertion) Name() string {
	return "Dash Insertion"
}

func (n *DashInsertion) Description() string {
	return "Inserting hyphens in the target domain"
}

func (n *DashInsertion) Fields() []string {
	return []string{}
}

func (n *DashInsertion) Headers() []string {
	return []string{}
}

func (n *DashInsertion) Exec(in urlinsane.Typo) (out []urlinsane.Typo) {
	out = append(out, in)
	return
}

// Register the plugin
func init() {
	algorithms.Add(CODE, func() urlinsane.Algorithm {
		return &DashInsertion{
			types: []string{algorithms.ENTITY, algorithms.DOMAIN},
		}
	})
}
