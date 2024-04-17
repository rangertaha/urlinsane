package none

import (
	typo "github.com/rangertaha/urlinsane"
	"github.com/rangertaha/urlinsane/plugins/algorithms"
)

type Homoglyphs struct {
	// Code() string
	// Name() string
	// Description() string
	// Fields() []string
	// Exec() func(Result) []Result
}

func (n *Homoglyphs) Code() string {
	return "hg"
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

func (n *Homoglyphs) Exec(typo.Result) (results []typo.Result) {
	return
}

// Register the plugin
func init() {
	algorithms.Add("hg", func() typo.Module {
		return &Homoglyphs{}
	})
}
