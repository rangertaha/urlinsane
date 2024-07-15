package AdjacentCharacterInsertion

import (
	typo "github.com/rangertaha/urlinsane"
	"github.com/rangertaha/urlinsane/plugins/algorithms"
)

type AdjacentCharacterInsertion struct {
	// Code() string
	// Name() string
	// Description() string
	// Fields() []string
	// Exec() func(Result) []Result
}

func (n *AdjacentCharacterInsertion) Code() string {
	return "aci"
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

func (n *AdjacentCharacterInsertion) Exec(typo.Result) (results []typo.Result) {
	return
}

// Register the plugin
func init() {
	algorithms.Add("aci", func() typo.Module {
		return &AdjacentCharacterInsertion{}
	})
}
