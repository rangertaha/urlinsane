package mds

import (
	"github.com/rangertaha/urlinsane"
	algorithms "github.com/rangertaha/urlinsane/plugins/algorithms"
)

type MissingDashes struct {
	types []string
}

func (n *MissingDashes) Code() string {
	return "mds"
}
func (n *MissingDashes) IsType(str string) bool {
	return algorithms.IsType(n.types, str)
}

func (n *MissingDashes) Name() string {
	return "Missing Dashes"
}

func (n *MissingDashes) Description() string {
	return "created by stripping all dashes from the domain"
}

func (n *MissingDashes) Fields() []string {
	return []string{}
}

func (n *MissingDashes) Headers() []string {
	return []string{}
}

func (n *MissingDashes) Exec(urlinsane.Typo) (results []urlinsane.Typo) {
	return
}

// Register the plugin
func init() {
	algorithms.Add("mds", func() urlinsane.Algorithm {
		return &MissingDashes{
			types: []string{algorithms.ENTITY, algorithms.DOMAINS},
		}
	})
}
