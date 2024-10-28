package mds
// Missing Dashes typos are created by omitting a dash from the domain.
// For example, www.a-b-c.com becomes www.ab-c.com, www.a-bc.com, and ww.abc.com



import (
	"github.com/rangertaha/urlinsane"
	"github.com/rangertaha/urlinsane/plugins/algorithms"
	"github.com/rangertaha/urlinsane/utils/nlp"
)

const CODE = "mds"

type MissingDashes struct {
	types []string
}

func (n *MissingDashes) Id() string {
	return CODE
}
func (n *MissingDashes) IsType(str string) bool {
	return algorithms.IsType(n.types, str)
}

func (n *MissingDashes) Name() string {
	return "Missing Dashes"
}

func (n *MissingDashes) Description() string {
	return "created by stripping all dashes from the name"
}

func (n *MissingDashes) Exec(typo urlinsane.Typo) (typos []urlinsane.Typo) {
	for _, variant := range nlp.MissingCharFunc(typo.Original().Repr(), "-") {
		if typo.Original().Repr() != variant {
			typos = append(typos, typo.New(variant))
		}
	}
	return
}

// Register the plugin
func init() {
	algorithms.Add(CODE, func() urlinsane.Algorithm {
		return &MissingDashes{
			types: []string{algorithms.ENTITY, algorithms.DOMAIN},
		}
	})
}

// // missingDashFunc typos are created by omitting a dash from the domain.
// // For example, www.a-b-c.com becomes www.ab-c.com, www.a-bc.com, and ww.abc.com
// func missingDashFunc(tc Result) (results []Result) {
// 	for _, str := range missingCharFunc(tc.Original.Domain, "-") {
// 		if tc.Original.Domain != str {
// 			dm := Domain{tc.Original.Subdomain, str, tc.Original.Suffix, Meta{}, false}
// 			results = append(results, Result{Original: tc.Original, Variant: dm, Typo: tc.Typo, Data: tc.Data})
// 		}
// 	}
// 	dm := Domain{tc.Original.Subdomain, strings.Replace(tc.Original.Domain, "-", "", -1), tc.Original.Suffix, Meta{}, false}
// 	results = append(results, Result{Original: tc.Original, Variant: dm, Typo: tc.Typo, Data: tc.Data})
// 	return results
// }
