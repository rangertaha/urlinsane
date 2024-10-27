package hp

import (
	"github.com/rangertaha/urlinsane"
	"github.com/rangertaha/urlinsane/plugins/algorithms"
)

const CODE = "hp"

type Homophones struct {
	types []string
}

func (n *Homophones) Code() string {
	return CODE
}
func (n *Homophones) IsType(str string) bool {
	return algorithms.IsType(n.types, str)
}

func (n *Homophones) Name() string {
	return "Homophones"
}

func (n *Homophones) Description() string {
	return "Created from sets of words that sound the same"
}

func (n *Homophones) Fields() []string {
	return []string{}
}

func (n *Homophones) Headers() []string {
	return []string{}
}

func (n *Homophones) Exec(in urlinsane.Typo) (out []urlinsane.Typo) {
	out = append(out, in)
	return
}

// Register the plugin
func init() {
	algorithms.Add(CODE, func() urlinsane.Algorithm {
		return &Homophones{
			types: []string{algorithms.ENTITY, algorithms.DOMAIN},
		}
	})
}
