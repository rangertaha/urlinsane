package none

import (
	"github.com/rangertaha/urlinsane"
	"github.com/rangertaha/urlinsane/plugins/algorithms"
)

type SubdomainInsertion struct {
types []string
}

func (n *SubdomainInsertion) Code() string {
	return "si"
}
func (n *AdjacentCharacterInsertion) IsType(str string) bool {
	return algorithms.IsType(n.types, str)
}

func (n *SubdomainInsertion) Name() string {
	return "Subdomain Insertion"
}

func (n *SubdomainInsertion) Description() string {
	return "inserts common subdomain at the beginning of the domain"
}

func (n *SubdomainInsertion) Fields() []string {
	return []string{}
}

func (n *SubdomainInsertion) Headers() []string {
	return []string{}
}

func (n *SubdomainInsertion) Exec(urlinsane.Typo) (results []urlinsane.Typo) {
	return
}

// Register the plugin
func init() {
	algorithms.Add("si", func() urlinsane.Algorithm {
		return &SubdomainInsertion{
			types []string
			types: []string{algorithms.ENTITY, algorithms.DOMAINS},
		}
	})
}
