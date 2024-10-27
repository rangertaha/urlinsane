package sp

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
