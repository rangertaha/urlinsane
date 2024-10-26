package dcr

import (
	"github.com/rangertaha/urlinsane"
	"github.com/rangertaha/urlinsane/plugins/algorithms"
)

type DoubleCharacterReplacement struct {
types []string
}

func (n *DoubleCharacterReplacement) Code() string {
	return "dcr"
}
func (n *DoubleCharacterReplacement) IsType(str string) bool {
	return algorithms.IsType(n.types, str)
}

func (n *DoubleCharacterReplacement) Name() string {
	return "DoubleCharacterReplacement"
}

func (n *DoubleCharacterReplacement) Description() string {
	return "Double Character Replacement repeats a character twice"
}

func (n *DoubleCharacterReplacement) Fields() []string {
	return []string{}
}

func (n *DoubleCharacterReplacement) Headers() []string {
	return []string{}
}

func (n *DoubleCharacterReplacement) Exec(urlinsane.Typo) (results []urlinsane.Typo) {
	return
}

// Register the plugin
func init() {
	algorithms.Add("dcr", func() urlinsane.Algorithm {
		return &DoubleCharacterReplacement{
			types: []string{algorithms.ENTITY, algorithms.DOMAINS},
		}
	})
}
