package none

import (
	typo "github.com/rangertaha/urlinsane"
	"github.com/rangertaha/urlinsane/plugins/algorithms"
)

type BitFlipping struct {
	// Code() string
	// Name() string
	// Description() string
	// Fields() []string
	// Exec() func(Result) []Result
}

func (n *BitFlipping) Code() string {
	return "bf"
}

func (n *BitFlipping) Name() string {
	return "Bit Flipping"
}

func (n *BitFlipping) Description() string {
	return "Relies on random bit-errors to redirect connections"
}

func (n *BitFlipping) Fields() []string {
	return []string{}
}

func (n *BitFlipping) Headers() []string {
	return []string{}
}

func (n *BitFlipping) Exec(typo.Result) (results []typo.Result) {
	return
}

// Register the plugin
func init() {
	algorithms.Add("bf", func() typo.Module {
		return &BitFlipping{}
	})
}
