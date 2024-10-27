package sp

import (
	"github.com/rangertaha/urlinsane"
	"github.com/rangertaha/urlinsane/plugins/algorithms"
)

const CODE = "sp"

type SingularPluralize struct {
	types []string
}

func (n *SingularPluralize) Code() string {
	return CODE
}
func (n *SingularPluralize) IsType(str string) bool {
	return algorithms.IsType(n.types, str)
}

func (n *SingularPluralize) Name() string {
	return "Singular Pluralize"
}

func (n *SingularPluralize) Description() string {
	return "Singular Pluralise creates a singular domain plural and vice versa"
}

func (n *SingularPluralize) Fields() []string {
	return []string{}
}

func (n *SingularPluralize) Headers() []string {
	return []string{}
}

func (n *SingularPluralize) Exec(urlinsane.Typo) (results []urlinsane.Typo) {
	return
}

// Register the plugin
func init() {
	algorithms.Add(CODE, func() urlinsane.Algorithm {
		return &SingularPluralize{
			types: []string{algorithms.ENTITY, algorithms.DOMAINS},
		}
	})
}
