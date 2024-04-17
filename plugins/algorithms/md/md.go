package none

import (
	typo "github.com/rangertaha/urlinsane"
	algorithms "github.com/rangertaha/urlinsane/plugins/algorithms"
)

type MissingDot struct {
	// Code() string
	// Name() string
	// Description() string
	// Fields() []string
	// Exec() func(Result) []Result
}

func (n *MissingDot) Code() string {
	return "md"
}

func (n *MissingDot) Name() string {
	return "Missing Dot"
}

func (n *MissingDot) Description() string {
	return "Missing Dot is created by omitting a dot from the domain"
}

func (n *MissingDot) Fields() []string {
	return []string{}
}

func (n *MissingDot) Headers() []string {
	return []string{}
}

func (n *MissingDot) Exec(typo.Result) (results []typo.Result) {
	return
}

// Register the plugin
func init() {
	algorithms.Add("md", func() typo.Module {
		return &MissingDot{}
	})
}
