package w2tld

import (
	"github.com/rangertaha/urlinsane"
	"github.com/rangertaha/urlinsane/plugins/algorithms"
)

const CODE = "w2tld"

type WrongThirdTLD struct {
	types []string
}

func (n *WrongThirdTLD) Code() string {
	return CODE
}
func (n *WrongThirdTLD) IsType(str string) bool {
	return algorithms.IsType(n.types, str)
}

func (n *WrongThirdTLD) Name() string {
	return "Wrong 2nd TLD"
}

func (n *WrongThirdTLD) Description() string {
	return "Wrong Second Level Domain"
}

func (n *WrongThirdTLD) Fields() []string {
	return []string{}
}

func (n *WrongThirdTLD) Headers() []string {
	return []string{}
}

func (n *WrongThirdTLD) Exec(urlinsane.Typo) (results []urlinsane.Typo) {
	return
}

// Register the plugin
func init() {
	algorithms.Add(CODE, func() urlinsane.Algorithm {
		return &WrongThirdTLD{
			types: []string{algorithms.ENTITY, algorithms.DOMAINS},
		}
	})
}
