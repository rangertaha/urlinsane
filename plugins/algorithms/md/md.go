package none

import (
	"github.com/rangertaha/urlinsane"
	algorithms "github.com/rangertaha/urlinsane/plugins/algorithms"
)

type MissingDot struct {
types []string
}

func (n *MissingDot) Code() string {
	return "md"
}
func (n *MissingDot) IsType(str string) bool {
	return algorithms.IsType(n.types, str)
}

func (n *MissingDot) Name() string {
	return "Missing Dot"
}

func (n *MissingDot) Description() string {
	return "Missing Dot is created by omitting a dot from the domain"
}

func (n *MissingDot) Fields() []string {
	return []string{}
}

func (n *MissingDot) Headers() []string {
	return []string{}
}

func (n *MissingDot) Exec(urlinsane.Typo) (results []urlinsane.Typo) {
	return
}

// Register the plugin
func init() {
	algorithms.Add("md", func() urlinsane.Algorithm {
		return &MissingDot{
			types: []string{algorithms.ENTITY, algorithms.DOMAINS},
		}
	})
}
