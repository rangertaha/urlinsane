package none

import (
	typo "github.com/rangertaha/urlinsane"
	"github.com/rangertaha/urlinsane/plugins/algorithms"
)

type VowelSwapping struct {
	// Code() string
	// Name() string
	// Description() string
	// Fields() []string
	// Exec() func(Result) []Result
}

func (n *VowelSwapping) Code() string {
	return "vs"
}

func (n *VowelSwapping) Name() string {
	return "Vowel Swapping"
}

func (n *VowelSwapping) Description() string {
	return "Vowel Swapping is created by swaps vowels"
}

func (n *VowelSwapping) Fields() []string {
	return []string{}
}

func (n *VowelSwapping) Exec(typo.Result) (results []typo.Result) {
	return
}

// Register the plugin
func init() {
	algorithms.Add("vs", func() typo.Module {
		return &VowelSwapping{

		}
	})
}
