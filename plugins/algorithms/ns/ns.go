package ns

import (
	"github.com/rangertaha/urlinsane"
	"github.com/rangertaha/urlinsane/plugins/algorithms"
)

type NumeralSwap struct {
types []string
}

func (n *NumeralSwap) Code() string {
	return "ns"
}
func (n *NumeralSwap) IsType(str string) bool {
	return algorithms.IsType(n.types, str)
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

func (n *NumeralSwap) Headers() []string {
	return []string{}
}

func (n *NumeralSwap) Exec(urlinsane.Typo) (results []urlinsane.Typo) {
	return
}

// Register the plugin
func init() {
	algorithms.Add("ns", func() urlinsane.Algorithm {
		return &NumeralSwap{
			types: []string{algorithms.ENTITY, algorithms.DOMAINS},
		}
	})
}
