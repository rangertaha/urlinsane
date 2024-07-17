package none

import (
	typo "github.com/rangertaha/urlinsane"
	"github.com/rangertaha/urlinsane/plugins/algorithms"
)

type CommonMisspellings struct {
	// Code() string
	// Name() string
	// Description() string
	// Fields() []string
	// Exec() func(Result) []Result
}

func (n *CommonMisspellings) Code() string {
	return "cm"
}

func (n *CommonMisspellings) Name() string {
	return "Common Misspellings"
}

func (n *CommonMisspellings) Description() string {
	return "Common Misspellings are created from a dictionary of commonly misspelled words"
}

func (n *CommonMisspellings) Fields() []string {
	return []string{}
}

func (n *CommonMisspellings) Headers() []string {
	return []string{}
}

func (n *CommonMisspellings) Exec(typo.Result) (results []typo.Result) {
	return
}

// Register the plugin
func init() {
	algorithms.Add("cm", func() typo.Module {
		return &CommonMisspellings{}
	})
}
