package none

import (
	"github.com/rangertaha/urlinsane"
	"github.com/rangertaha/urlinsane/plugins/algorithms"
)

type VowelSwapping struct {
types []string
}

func (n *VowelSwapping) Code() string {
	return "vs"
}
func (n *AdjacentCharacterInsertion) IsType(str string) bool {
	return algorithms.IsType(n.types, str)
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

func (n *VowelSwapping) Headers() []string {
	return []string{}
}

func (n *VowelSwapping) Exec(urlinsane.Typo) (results []urlinsane.Typo) {
	return
}

// Register the plugin
func init() {
	algorithms.Add("vs", func() urlinsane.Algorithm {
		return &VowelSwapping{
			types []string
			types: []string{algorithms.ENTITY, algorithms.DOMAINS},
		}
	})
}
