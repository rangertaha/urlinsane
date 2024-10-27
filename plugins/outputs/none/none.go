package none

import (
	"github.com/rangertaha/urlinsane"
	"github.com/rangertaha/urlinsane/plugins/information"
)

const CODE = "aci"

type None struct {
	types []string
}

func (n *None) Code() string {
	return CODE
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

func (n *None) Exec(urlinsane.Typo) (results urlinsane.Typo) {
	return
}

// Register the plugin
func init() {
	information.Add(CODE, func() urlinsane.Information {
		return &None{
			types: []string{information.ENTITY, information.DOMAINS},
		}
	})
}
