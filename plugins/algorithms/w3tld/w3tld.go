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
	return "w3tld"
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
		return &WrongThirdTLD{}
	})
}
