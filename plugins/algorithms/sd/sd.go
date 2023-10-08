package none

import (
	typo "github.com/rangertaha/urlinsane"
	"github.com/rangertaha/urlinsane/plugins/algorithms"
)

type StripDash struct {
	// Code() string
	// Name() string
	// Description() string
	// Fields() []string
	// Exec() func(Result) []Result
}

func (n *StripDash) Code() string {
	return "sd"
}

func (n *StripDash) Name() string {
	return "Strip Dash"
}

func (n *StripDash) Description() string {
	return "created by omitting a single dash from the domain"
}

func (n *StripDash) Fields() []string {
	return []string{}
}

func (n *StripDash) Headers() []string {
	return []string{}
}

func (n *StripDash) Exec(typo.Result) (results []typo.Result) {
	return
}

// Register the plugin
func init() {
	algorithms.Add("sd", func() typo.Module {
		return &StripDash{}
	})
}
