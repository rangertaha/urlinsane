package none

import (
	"github.com/rangertaha/urlinsane"
	"github.com/rangertaha/urlinsane/plugins/algorithms"
)

type WrongThirdTLD struct {
types []string
}

func (n *WrongThirdTLD) Code() string {
	return "w3tld"
}
func (n *WrongThirdTLD) IsType(str string) bool {
	return algorithms.IsType(n.types, str)
}

func (n *WrongThirdTLD) Name() string {
	return "Wrong 3rd TLD"
}

func (n *WrongThirdTLD) Description() string {
	return "Wrong Third Level Domain"
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
	algorithms.Add("w3tld", func() urlinsane.Algorithm {
		return &WrongThirdTLD{
			types []string
			types: []string{algorithms.ENTITY, algorithms.DOMAINS},
		}
	})
}
