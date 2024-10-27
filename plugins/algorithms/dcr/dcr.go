package dcr

import (
	"github.com/rangertaha/urlinsane"
	"github.com/rangertaha/urlinsane/plugins/algorithms"
)

const CODE = "dcr"

type DoubleCharacterReplacement struct {
	types []string
}

func (n *DoubleCharacterReplacement) Id() string {
	return CODE
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

func (n *DoubleCharacterReplacement) Exec(in urlinsane.Typo) (out []urlinsane.Typo) {
	out = append(out, in)
	return
}

// Register the plugin
func init() {
	algorithms.Add(CODE, func() urlinsane.Algorithm {
		return &DoubleCharacterReplacement{
			types: []string{algorithms.ENTITY, algorithms.DOMAIN},
		}
	})
}
