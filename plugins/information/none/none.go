package none

import (
	urli "github.com/rangertaha/urlinsane"
	"github.com/rangertaha/urlinsane/plugins/information"
)

type None struct {
	// Code() string
	// Name() string
	// Description() string
	// Fields() []string
	// Exec() func(Result) []Result
}

func (n *None) Code() string {
	return "none"
}

func (n *None) Name() string {
	return "None"
}

func (n *None) Description() string {
	return "---------------------------------"
}

func (n *None) Fields() []string {
	return []string{}
}

func (n *None) Headers() []string {
	return []string{}
}

func (n *None) Exec(urli.Result) (results []urli.Result) {
	return
}

// Register the plugin
func init() {
	information.Add("none", func() urli.Module {
		return &None{}
	})
}
