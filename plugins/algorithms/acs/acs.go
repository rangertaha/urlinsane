package AdjacentCharacterSubstitution

import (
	"github.com/rangertaha/urlinsane"
	typo "github.com/rangertaha/urlinsane"
	"github.com/rangertaha/urlinsane/plugins/algorithms"
)

type AdjacentCharacterSubstitution struct {
	// Code() string
	// Name() string
	// Description() string
	// Fields() []string
	// Exec() func(Result) []Result
}

func (n *AdjacentCharacterSubstitution) Code() string {
	return "acs"
}

func (n *AdjacentCharacterSubstitution) Name() string {
	return "Adjacent Character Substitution"
}

func (n *AdjacentCharacterSubstitution) Description() string {
	return ""
}

func (n *AdjacentCharacterSubstitution) Fields() []string {
	return []string{}
}

func (n *AdjacentCharacterSubstitution) Headers() []string {
	return []string{}
}

func (n *AdjacentCharacterSubstitution) Exec(urlinsane.Typo) (results []urlinsane.Typo) {
	return
}

// Register the plugin
func init() {
	algorithms.Add("acs", func() typo.Module {
		return &AdjacentCharacterSubstitution{}
	})
}
