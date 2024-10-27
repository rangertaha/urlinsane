package cm

import (
	"github.com/rangertaha/urlinsane"
	"github.com/rangertaha/urlinsane/plugins/algorithms"
)

const CODE = "cm"

type CommonMisspellings struct {
	types []string
}

func (n *CommonMisspellings) Code() string {
	return CODE
}
func (n *CommonMisspellings) IsType(str string) bool {
	return algorithms.IsType(n.types, str)
}

func (n *CommonMisspellings) Name() string {
	return "Common Misspellings"
}

func (n *CommonMisspellings) Description() string {
	return "Common Misspellings are created from a dictionary of commonly misspelled words"
}

func (n *CommonMisspellings) Fields() []string {
	return []string{}
}

func (n *CommonMisspellings) Headers() []string {
	return []string{}
}

func (n *CommonMisspellings) Exec(urlinsane.Typo) (results []urlinsane.Typo) {
	return
}

// Register the plugin
func init() {
	algorithms.Add(CODE, func() urlinsane.Algorithm {
		return &CommonMisspellings{
			types: []string{algorithms.ENTITY, algorithms.DOMAINS},
		}
	})
}
