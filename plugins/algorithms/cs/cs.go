package none

import (
	"github.com/rangertaha/urlinsane"
	typo "github.com/rangertaha/urlinsane"
	"github.com/rangertaha/urlinsane/plugins/algorithms"
)

type CharacterSwap struct {
	// Code() string
	// Name() string
	// Description() string
	// Fields() []string
	// Exec() func(Result) []Result
}

func (n *CharacterSwap) Code() string {
	return "cs"
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

func (n *CharacterSwap) Exec(urlinsane.Typo) (results []urlinsane.Typo) {
	return
}

// Register the plugin
func init() {
	algorithms.Add("cs", func() typo.Module {
		return &CharacterSwap{}
	})
}
