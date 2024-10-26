package none

import (
	"github.com/rangertaha/urlinsane"
	"github.com/rangertaha/urlinsane/plugins/algorithms"
)

type WrongThirdTLD struct {
	// Code() string
	// Name() string
	// Description() string
	// Fields() []string
	// Exec() func(Result) []Result
}

func (n *WrongThirdTLD) Code() string {
	return "w2tld"
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
	algorithms.Add("w2tld", func() urlinsane.Algorithm {
		return &WrongThirdTLD{}
	})
}
