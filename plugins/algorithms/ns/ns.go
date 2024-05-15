package none

import (
	typo "github.com/rangertaha/urlinsane"
	"github.com/rangertaha/urlinsane/plugins/algorithms"
)

type NumeralSwap struct {
	// Code() string
	// Name() string
	// Description() string
	// Fields() []string
	// Exec() func(Result) []Result
}

func (n *NumeralSwap) Code() string {
	return "ns"
}

func (n *NumeralSwap) Name() string {
	return "NumeralSwap"
}

func (n *NumeralSwap) Description() string {
	return "Numeral Swap numbers, words and vice versa"
}

func (n *NumeralSwap) Fields() []string {
	return []string{}
}

func (n *NumeralSwap) Exec(typo.Result) (results []typo.Result) {
	return
}

// Register the plugin
func init() {
	algorithms.Add("ns", func() typo.Module {
		return &NumeralSwap{

		}
	})
}
