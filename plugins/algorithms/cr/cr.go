package none

import (
	"github.com/rangertaha/urlinsane"
	"github.com/rangertaha/urlinsane/plugins/algorithms"
)

type CharacterRepeat struct {
types []string
}

func (n *CharacterRepeat) Code() string {
	return "cr"
}
func (n *CharacterRepeat) IsType(str string) bool {
	return algorithms.IsType(n.types, str)
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

func (n *CharacterRepeat) Exec(urlinsane.Typo) (results []urlinsane.Typo) {
	return
}

// Register the plugin
func init() {
	algorithms.Add("cr", func() urlinsane.Algorithm {
		return &CharacterRepeat{
			types: []string{algorithms.ENTITY, algorithms.DOMAINS},
		}
	})
}
