package none

import (
	"github.com/rangertaha/urlinsane"
	"github.com/rangertaha/urlinsane/plugins/algorithms"
)

type Homoglyphs struct {
	types []string
}

func (n *Homoglyphs) Code() string {
	return "hg"
}
func (n *Homoglyphs) IsType(str string) bool {
	return algorithms.IsType(n.types, str)
}

func (n *Homoglyphs) Name() string {
	return "Homoglyphs"
}

func (n *Homoglyphs) Description() string {
	return "Replaces characters with characters that look similar"
}

func (n *Homoglyphs) Fields() []string {
	return []string{}
}

func (n *Homoglyphs) Headers() []string {
	return []string{}
}

func (n *Homoglyphs) Exec(urlinsane.Typo) (results []urlinsane.Typo) {
	return
}

// Register the plugin
func init() {
	algorithms.Add("hg", func() urlinsane.Algorithm {
		return &Homoglyphs{
			types: []string{algorithms.ENTITY, algorithms.DOMAINS},
		}
	})
}
