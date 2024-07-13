package none

import (
	typo "github.com/rangertaha/urlinsane"
	"github.com/rangertaha/urlinsane/plugins/algorithms"
)

type None struct {
	// Code() string
	// Name() string
	// Description() string
	// Fields() []string
	// Exec() func(Result) []Result
}

func (n *None) Code() string {
	return "SD"
}

func (n *None) Name() string {
	return "Strip Dash"
}

func (n *None) Description() string {
	return "created by omitting a single dash from the domain"
}

func (n *None) Fields() []string {
	return []string{}
}

func (n *None) Exec(typo.Result) (results []typo.Result) {
	return
}

// Register the plugin
func init() {
	algorithms.Add("SD", func() typo.Module {
		return &None{

		}
	})
}
