package si

import (
	"github.com/rangertaha/urlinsane"
	"github.com/rangertaha/urlinsane/plugins/algorithms"
)

const CODE = "si"

type SubdomainInsertion struct {
	types []string
}

func (n *SubdomainInsertion) Id() string {
	return CODE
}
func (n *SubdomainInsertion) IsType(str string) bool {
	return algorithms.IsType(n.types, str)
}

func (n *SubdomainInsertion) Name() string {
	return "Subdomain Insertion"
}

func (n *SubdomainInsertion) Description() string {
	return "Inserts common subdomain at the beginning of the domain"
}

func (n *SubdomainInsertion) Fields() []string {
	return []string{}
}

func (n *SubdomainInsertion) Headers() []string {
	return []string{}
}

func (n *SubdomainInsertion) Exec(in urlinsane.Typo) (out []urlinsane.Typo) {
	out = append(out, in)
	return
}

// Register the plugin
func init() {
	algorithms.Add(CODE, func() urlinsane.Algorithm {
		return &SubdomainInsertion{
			types: []string{algorithms.DOMAIN},
		}
	})
}


// // subdomainInsertionFunc typos are created by inserting common subdomains at the begining of the domain. wwwgoogle.com and ftpgoogle.com
// func subdomainInsertionFunc(tc Result) (results []Result) {
// 	for _, str := range datasets.SUBDOMAINS {
// 		dm := Domain{tc.Original.Subdomain, str + tc.Original.Domain, tc.Original.Suffix, Meta{}, false}
// 		results = append(results, Result{Original: tc.Original, Variant: dm, Typo: tc.Typo, Data: tc.Data})
// 	}
// 	return results
// }
