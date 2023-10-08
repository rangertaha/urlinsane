package none

import (
	typo "github.com/rangertaha/urlinsane"
	"github.com/rangertaha/urlinsane/plugins/algorithms"
)

type DoubleCharacterReplacement struct {
	// Code() string
	// Name() string
	// Description() string
	// Fields() []string
	// Exec() func(Result) []Result
}

func (n *DoubleCharacterReplacement) Code() string {
	return "dcr"
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

func (n *DoubleCharacterReplacement) Exec(typo.Result) (results []typo.Result) {
	return
}

// Register the plugin
func init() {
	algorithms.Add("dcr", func() typo.Module {
		return &DoubleCharacterReplacement{}
	})
}
