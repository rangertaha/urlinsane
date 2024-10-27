package ia

import (
	"github.com/rangertaha/urlinsane"
	"github.com/rangertaha/urlinsane/plugins/algorithms"
)

const CODE = "ia"

type AlphabetInsertion struct {
	types []string
}

func (n *AlphabetInsertion) Id() string {
	return CODE
}
func (n *AlphabetInsertion) IsType(str string) bool {
	return algorithms.IsType(n.types, str)
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

func (n *AlphabetInsertion) Exec(in urlinsane.Typo) (out []urlinsane.Typo) {
	out = append(out, in)
	return
}

// Register the plugin
func init() {
	algorithms.Add(CODE, func() urlinsane.Algorithm {
		return &AlphabetInsertion{
			types: []string{algorithms.ENTITY, algorithms.DOMAIN},
		}
	})
}
