package aci

import (
	"github.com/rangertaha/urlinsane"
	"github.com/rangertaha/urlinsane/plugins/algorithms"
)

const CODE = "aci"

type AdjacentCharacterInsertion struct {
	types []string
}

func (n *AdjacentCharacterInsertion) Id() string {
	return CODE
}

func (n *AdjacentCharacterInsertion) IsType(str string) bool {
	return algorithms.IsType(n.types, str)
}

func (n *AdjacentCharacterInsertion) Name() string {
	return "Adjacent Character Insertion"
}

func (n *AdjacentCharacterInsertion) Description() string {
	return "Adjacent Character Insertion inserts adjacent character"
}

func (n *AdjacentCharacterInsertion) Fields() []string {
	return []string{}
}

func (n *AdjacentCharacterInsertion) Headers() []string {
	return []string{}
}

func (n *AdjacentCharacterInsertion) Exec(in urlinsane.Typo) (out []urlinsane.Typo) {
	out = append(out, in)
	return
}

// Register the plugin
func init() {
	algorithms.Add(CODE, func() urlinsane.Algorithm {
		return &AdjacentCharacterInsertion{
			types: []string{algorithms.ENTITY, algorithms.DOMAIN},
		}
	})
}
