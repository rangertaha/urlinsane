package none

import (
	"github.com/rangertaha/urlinsane"
	"github.com/rangertaha/urlinsane/plugins/algorithms"
)

type Homophones struct {
	// Code() string
	// Name() string
	// Description() string
	// Fields() []string
	// Exec() func(Result) []Result
}

func (n *Homophones) Code() string {
	return "hp"
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

func (n *Homophones) Exec(urlinsane.Typo) (results []urlinsane.Typo) {
	return
}

// Register the plugin
func init() {
	algorithms.Add("hp", func() urlinsane.Algorithm {
		return &Homophones{}
	})
}
