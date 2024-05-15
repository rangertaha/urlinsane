package none

import (
	typo "github.com/rangertaha/urlinsane"
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

func (n *WrongThirdTLD) Exec(typo.Result) (results []typo.Result) {
	return
}

// Register the plugin
func init() {
	algorithms.Add("w3tld", func() typo.Module {
		return &WrongThirdTLD{

		}
	})
}
