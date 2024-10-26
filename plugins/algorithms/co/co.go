package none

import (
	"github.com/rangertaha/urlinsane"
	typo "github.com/rangertaha/urlinsane"
	"github.com/rangertaha/urlinsane/plugins/algorithms"
)

type CharacterOmission struct {
	// Code() string
	// Name() string
	// Description() string
	// Fields() []string
	// Exec() func(Result) []Result
}

func (n *CharacterOmission) Code() string {
	return "co"
}

func (n *CharacterOmission) Name() string {
	return "Character Omission"
}

func (n *CharacterOmission) Description() string {
	return "omitting a character from the domain"
}

func (n *CharacterOmission) Fields() []string {
	return []string{}
}

func (n *CharacterOmission) Headers() []string {
	return []string{}
}

func (n *CharacterOmission) Exec(urlinsane.Typo) (results []urlinsane.Typo) {
	return
}

// Register the plugin
func init() {
	algorithms.Add("co", func() typo.Module {
		return &CharacterOmission{}
	})
}
