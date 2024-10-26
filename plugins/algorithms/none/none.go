package none

import (
	"github.com/rangertaha/urlinsane"
	"github.com/rangertaha/urlinsane/plugins/algorithms"
)

type None struct {
	types []string
}

func (n *None) Code() string {
	return ""
}
func (n *None) IsType(str string) bool {
	return algorithms.IsType(n.types, str)
}

func (n *None) Name() string {
	return "None"
}

func (n *None) Description() string {
	return ""
}

func (n *None) Fields() []string {
	return []string{}
}

func (n *None) Headers() []string {
	return []string{}
}

func (n *None) Exec(urlinsane.Typo) (results []urlinsane.Typo) {
	return
}

// Register the plugin
func init() {
	algorithms.Add("none", func() urlinsane.Algorithm {
		return &None{
			types: []string{algorithms.ENTITY, algorithms.DOMAINS},
		}
	})
}
