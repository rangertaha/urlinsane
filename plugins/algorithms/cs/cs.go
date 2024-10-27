package cs

import (
	"github.com/rangertaha/urlinsane"
	"github.com/rangertaha/urlinsane/plugins/algorithms"
)

const CODE = "cs"

type CharacterSwap struct {
	types []string
}

func (n *CharacterSwap) Code() string {
	return CODE
}
func (n *CharacterSwap) IsType(str string) bool {
	return algorithms.IsType(n.types, str)
}

func (n *CharacterSwap) Name() string {
	return "Character Swap"
}

func (n *CharacterSwap) Description() string {
	return "Character Swap Swapping two consecutive characters in a domain"
}

func (n *CharacterSwap) Fields() []string {
	return []string{}
}

func (n *CharacterSwap) Headers() []string {
	return []string{}
}

func (n *CharacterSwap) Exec(in urlinsane.Typo) (out []urlinsane.Typo) {
	out = append(out, in)
	return
}

// Register the plugin
func init() {
	algorithms.Add(CODE, func() urlinsane.Algorithm {
		return &CharacterSwap{
			types: []string{algorithms.ENTITY, algorithms.DOMAIN},
		}
	})
}
