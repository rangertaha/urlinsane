package AlphabetReplacement

import (
	"github.com/rangertaha/urlinsane"
	"github.com/rangertaha/urlinsane/plugins/algorithms"
)

type AlphabetReplacement struct {
types []string
}

func (n *AlphabetReplacement) Code() string {
	return "ar"
}
func (n *AlphabetReplacement) IsType(str string) bool {
	return algorithms.IsType(n.types, str)
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

func (n *AlphabetReplacement) Headers() []string {
	return []string{}
}

func (n *AlphabetReplacement) Exec(urlinsane.Typo) (results []urlinsane.Typo) {
	return
}

// Register the plugin
func init() {
	algorithms.Add("ar", func() urlinsane.Algorithm {
		return &AlphabetReplacement{
			
			types: []string{algorithms.ENTITY, algorithms.DOMAINS},
		}
	})
}
