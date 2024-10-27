package wtld

import (
	"github.com/rangertaha/urlinsane"
	"github.com/rangertaha/urlinsane/plugins/algorithms"
)

const CODE = "wtld"

type WrongTLD struct {
	types []string
}

func (n *WrongTLD) Code() string {
	return CODE
}
func (n *WrongTLD) IsType(str string) bool {
	return algorithms.IsType(n.types, str)
}
func (n *WrongTLD) Name() string {
	return "Wrong TLD"
}

func (n *WrongTLD) Description() string {
	return "Wrong top level domain (TLD)"
}

func (n *WrongTLD) Fields() []string {
	return []string{}
}

func (n *WrongTLD) Headers() []string {
	return []string{}
}

func (n *WrongTLD) Exec(in urlinsane.Typo) (out []urlinsane.Typo) {
	out = append(out, in)
	return
}

// Register the plugin
func init() {
	algorithms.Add(CODE, func() urlinsane.Algorithm {
		return &WrongTLD{
			types: []string{algorithms.DOMAIN},
		}
	})
}
