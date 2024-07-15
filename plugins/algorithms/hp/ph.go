package none

import (
	typo "github.com/rangertaha/urlinsane"
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

func (n *Homophones) Exec(typo.Result) (results []typo.Result) {
	return
}

// Register the plugin
func init() {
	algorithms.Add("hp", func() typo.Module {
		return &Homophones{

		}
	})
}
