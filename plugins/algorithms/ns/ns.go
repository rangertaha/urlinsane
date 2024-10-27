package ns

import (
	"github.com/rangertaha/urlinsane"
	"github.com/rangertaha/urlinsane/plugins/algorithms"
)

const CODE = "ns"

type NumeralSwap struct {
	types []string
}

func (n *NumeralSwap) Id() string {
	return CODE
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

func (n *NumeralSwap) Exec(in urlinsane.Typo) (out []urlinsane.Typo) {
	out = append(out, in)
	return
}

// Register the plugin
func init() {
	algorithms.Add(CODE, func() urlinsane.Algorithm {
		return &NumeralSwap{
			types: []string{algorithms.ENTITY, algorithms.DOMAIN},
		}
	})
}
