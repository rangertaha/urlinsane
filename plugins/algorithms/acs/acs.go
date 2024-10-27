package acs

import (
	"github.com/rangertaha/urlinsane"
	"github.com/rangertaha/urlinsane/plugins/algorithms"
)

const CODE = "acs"

type AdjacentCharacterSubstitution struct {
	types []string
}

func (n *AdjacentCharacterSubstitution) Id() string {
	return CODE
}

func (n *AdjacentCharacterSubstitution) IsType(str string) bool {
	return algorithms.IsType(n.types, str)
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

func (n *AdjacentCharacterSubstitution) Exec(in urlinsane.Typo) (out []urlinsane.Typo) {
	out = append(out, in)
	return
}

// Register the plugin
func init() {
	algorithms.Add(CODE, func() urlinsane.Algorithm {
		return &AdjacentCharacterSubstitution{
			types: []string{algorithms.ENTITY, algorithms.DOMAIN},
		}
	})
}
