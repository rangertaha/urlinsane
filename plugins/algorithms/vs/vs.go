package vs

import (
	"github.com/rangertaha/urlinsane"
	"github.com/rangertaha/urlinsane/plugins/algorithms"
)

const CODE = "vs"

type VowelSwapping struct {
	types []string
}

func (n *VowelSwapping) Code() string {
	return CODE
}
func (n *VowelSwapping) IsType(str string) bool {
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

func (n *VowelSwapping) Exec(in urlinsane.Typo) (out []urlinsane.Typo) {
	out = append(out, in)
	return
}

// Register the plugin
func init() {
	algorithms.Add(CODE, func() urlinsane.Algorithm {
		return &VowelSwapping{
			types: []string{algorithms.ENTITY, algorithms.DOMAIN},
		}
	})
}
