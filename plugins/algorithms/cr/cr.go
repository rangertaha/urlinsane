package none

import (
	typo "github.com/rangertaha/urlinsane"
	"github.com/rangertaha/urlinsane/plugins/algorithms"
)

type CharacterRepeat struct {
	// Code() string
	// Name() string
	// Description() string
	// Fields() []string
	// Exec() func(Result) []Result
}

func (n *CharacterRepeat) Code() string {
	return "cr"
}

func (n *CharacterRepeat) Name() string {
	return "CharacterRepeat"
}

func (n *CharacterRepeat) Description() string {
	return "Character Repeat Repeats a character of the domain name twice"
}

func (n *CharacterRepeat) Fields() []string {
	return []string{}
}

func (n *CharacterRepeat) Headers() []string {
	return []string{}
}

func (n *CharacterRepeat) Exec(typo.Result) (results []typo.Result) {
	return
}

// Register the plugin
func init() {
	algorithms.Add("cr", func() typo.Module {
		return &CharacterRepeat{}
	})
}
