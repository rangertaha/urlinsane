package AlphabetReplacement

import (
	typo "github.com/rangertaha/urlinsane"
	"github.com/rangertaha/urlinsane/plugins/algorithms"
)

type AlphabetReplacement struct {
	// Code() string
	// Name() string
	// Description() string
	// Fields() []string
	// Exec() func(Result) []Result
}

func (n *AlphabetReplacement) Code() string {
	return "ar"
}

func (n *AlphabetReplacement) Name() string {
	return "Alphabet Replacement"
}

func (n *AlphabetReplacement) Description() string {
	return "Replaces an alphabet in the target domain"
}

func (n *AlphabetReplacement) Fields() []string {
	return []string{}
}

func (n *AlphabetReplacement) Exec(typo.Result) (results []typo.Result) {
	return
}

// Register the plugin
func init() {
	algorithms.Add("ar", func() typo.Module {
		return &AlphabetReplacement{

		}
	})
}
