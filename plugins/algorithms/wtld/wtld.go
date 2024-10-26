package none

import (
	"github.com/rangertaha/urlinsane"
	"github.com/rangertaha/urlinsane/plugins/algorithms"
)

type WrongTLD struct {
types []string
}

func (n *WrongTLD) Code() string {
	return "wtld"
}
func (n *WrongThirdTLD) IsType(str string) bool {
	return algorithms.IsType(n.types, str)
}
func (n *WrongTLD) Name() string {
	return "Wrong TLD"
}

func (n *WrongTLD) Description() string {
	return "Wrong Top Level Domain"
}

func (n *WrongTLD) Fields() []string {
	return []string{}
}

func (n *WrongTLD) Headers() []string {
	return []string{}
}

func (n *WrongTLD) Exec(urlinsane.Typo) (results []urlinsane.Typo) {
	return
}

// Register the plugin
func init() {
	algorithms.Add("wtld", func() urlinsane.Algorithm {
		return &WrongTLD{
			types []string
			types: []string{algorithms.ENTITY, algorithms.DOMAINS},
		}
	})
}
