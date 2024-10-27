package idn

// An internationalized domain name (IDN) is a domain name that includes characters
// outside of the Latin alphabet, such as letters from Arabic, Chinese, Cyrillic, or
// Devanagari scripts. IDNs allow users to use domain names in their local languages
// and scripts.

import (
	"github.com/rangertaha/urlinsane"
	"github.com/rangertaha/urlinsane/plugins/algorithms"
)

const CODE = "idn"

// // Idna ...
// func (d *Domain) Idna() (punycode string) {
// 	punycode, _ = idna.Punycode.ToASCII(d.String())
// 	return
// }

type InternationalizedDomainName struct {
	types []string
}

func (n *InternationalizedDomainName) Code() string {
	return CODE
}
func (n *InternationalizedDomainName) IsType(str string) bool {
	return algorithms.IsType(n.types, str)
}

func (n *InternationalizedDomainName) Name() string {
	return "Internationalized Domain Name"
}

func (n *InternationalizedDomainName) Description() string {
	return "Internationalized domain names includes characters outside of the Latin alphabet"
}

func (n *InternationalizedDomainName) Fields() []string {
	return []string{}
}

func (n *InternationalizedDomainName) Headers() []string {
	return []string{}
}

func (n *InternationalizedDomainName) Exec(in urlinsane.Typo) (out []urlinsane.Typo) {
	out = append(out, in)
	return
}

// Register the plugin
func init() {
	algorithms.Add(CODE, func() urlinsane.Algorithm {
		return &InternationalizedDomainName{

			types: []string{algorithms.DOMAIN},
		}
	})
}
