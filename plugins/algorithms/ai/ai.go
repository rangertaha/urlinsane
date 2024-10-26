package AlphabetInsertion

import (
	"github.com/rangertaha/urlinsane"
	"github.com/rangertaha/urlinsane/plugins/algorithms"
)

type AlphabetInsertion struct {
	// Code() string
	// Name() string
	// Description() string
	// Fields() []string
	// Exec() func(Result) []Result
}

func (n *AlphabetInsertion) Code() string {
	return "ai"
}

func (n *AlphabetInsertion) Name() string {
	return "Alphabet Insertion"
}

func (n *AlphabetInsertion) Description() string {
	return "Inserting the language specific alphabet in the target domain"
}

func (n *AlphabetInsertion) Fields() []string {
	return []string{}
}

func (n *AlphabetInsertion) Headers() []string {
	return []string{}
}

func (n *AlphabetInsertion) Exec(urlinsane.Typo) (results []urlinsane.Typo) {
	return
}

// Register the plugin
func init() {
	algorithms.Add("AlphabetInsertion", func() urlinsane.Algorithm {
		return &AlphabetInsertion{}
	})
}
