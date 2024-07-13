package none

import (
	typo "github.com/rangertaha/urlinsane"
	algorithms "github.com/rangertaha/urlinsane/plugins/algorithms"
)

type None struct {
	// Code() string
	// Name() string
	// Description() string
	// Fields() []string
	// Exec() func(Result) []Result
}

func (n *None) Code() string {
	return "MDS"
}

func (n *None) Name() string {
	return "Missing Dashes"
}

func (n *None) Description() string {
	return "created by stripping all dashes from the domain"
}

func (n *None) Fields() []string {
	return []string{}
}

func (n *None) Exec(typo.Result) (results []typo.Result) {
	return
}

// Register the plugin
func init() {
	algorithms.Add("MDS", func() typo.Module {
		return &None{}
	})
}
