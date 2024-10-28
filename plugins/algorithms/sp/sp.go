package sp

// singularPluraliseFunc are created by making a singular domain plural and
// vice versa. For example, www.google.com becomes www.googles.com and
// www.games.co.nz becomes www.game.co.nz
// func singularPluraliseFunc(tc Result) (results []Result) {
// 	for _, pchar := range []string{"s", "ing"} {
// 		var domain string
// 		if strings.HasSuffix(tc.Original.Domain, pchar) {
// 			domain = strings.TrimSuffix(tc.Original.Domain, pchar)
// 		} else {
// 			domain = tc.Original.Domain + pchar
// 		}
// 		dm := Domain{tc.Original.Subdomain, domain, tc.Original.Suffix, Meta{}, false}
// 		results = append(results, Result{Original: tc.Original, Variant: dm, Typo: tc.Typo, Data: tc.Data})
// 	}
// 	return
// }

import (
	"github.com/rangertaha/urlinsane"
	"github.com/rangertaha/urlinsane/plugins/algorithms"
)

const CODE = "sp"

type SingularPluralize struct {
	types []string
}

func (n *SingularPluralize) Id() string {
	return CODE
}
func (n *SingularPluralize) IsType(str string) bool {
	return algorithms.IsType(n.types, str)
}

func (n *SingularPluralize) Name() string {
	return "Singular Pluralize"
}

func (n *SingularPluralize) Description() string {
	return "Creates singular and plural names"
}

func (n *SingularPluralize) Fields() []string {
	return []string{}
}

func (n *SingularPluralize) Headers() []string {
	return []string{}
}

func (n *SingularPluralize) Exec(in urlinsane.Typo) (out []urlinsane.Typo) {
	out = append(out, in)
	return
}

// Register the plugin
func init() {
	algorithms.Add(CODE, func() urlinsane.Algorithm {
		return &SingularPluralize{
			types: []string{algorithms.ENTITY, algorithms.DOMAIN},
		}
	})
}
