package mds

import (
	"github.com/rangertaha/urlinsane"
	"github.com/rangertaha/urlinsane/plugins/algorithms"
)

const CODE = "mds"

type MissingDashes struct {
	types []string
}

func (n *MissingDashes) Code() string {
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

func (n *MissingDashes) Fields() []string {
	return []string{}
}

func (n *MissingDashes) Headers() []string {
	return []string{}
}

func (n *MissingDashes) Exec(in urlinsane.Typo) (out []urlinsane.Typo) {
	out = append(out, in)
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
