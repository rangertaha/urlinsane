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
	return "md"
}

func (n *None) Name() string {
	return "Missing Dot"
}

func (n *None) Description() string {
	return "Missing Dot is created by omitting a dot from the domain"
}

func (n *None) Fields() []string {
	return []string{}
}

func (n *None) Exec(typo.Result) (results []typo.Result) {
	return
}

// Register the plugin
func init() {
	algorithms.Add("md", func() typo.Module {
		return &None{}
	})
}
