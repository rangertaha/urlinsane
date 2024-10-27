package w3tld

import (
	"github.com/rangertaha/urlinsane"
	"github.com/rangertaha/urlinsane/plugins/algorithms"
)

const CODE = "w3tld"

type WrongThirdTLD struct {
	types []string
}

func (n *WrongThirdTLD) Id() string {
	return CODE
}
func (n *WrongThirdTLD) IsType(str string) bool {
	return algorithms.IsType(n.types, str)
}

func (n *WrongThirdTLD) Name() string {
	return "Wrong TLD3"
}

func (n *WrongThirdTLD) Description() string {
	return "Wrong third level domain (TLD3)"
}

func (n *WrongThirdTLD) Fields() []string {
	return []string{}
}

func (n *WrongThirdTLD) Headers() []string {
	return []string{}
}

func (n *WrongThirdTLD) Exec(in urlinsane.Typo) (out []urlinsane.Typo) {
	out = append(out, in)
	return
}

// Register the plugin
func init() {
	algorithms.Add(CODE, func() urlinsane.Algorithm {
		return &WrongThirdTLD{
			types: []string{algorithms.DOMAIN},
		}
	})
}
