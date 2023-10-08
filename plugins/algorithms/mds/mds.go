package none

import (
	typo "github.com/rangertaha/urlinsane"
	algorithms "github.com/rangertaha/urlinsane/plugins/algorithms"
)

type MissingDashes struct {
	// Code() string
	// Name() string
	// Description() string
	// Fields() []string
	// Exec() func(Result) []Result
}

func (n *MissingDashes) Code() string {
	return "mds"
}

func (n *MissingDashes) Name() string {
	return "Missing Dashes"
}

func (n *MissingDashes) Description() string {
	return "created by stripping all dashes from the domain"
}

func (n *MissingDashes) Fields() []string {
	return []string{}
}

func (n *MissingDashes) Headers() []string {
	return []string{}
}

func (n *MissingDashes) Exec(typo.Result) (results []typo.Result) {
	return
}

// Register the plugin
func init() {
	algorithms.Add("mds", func() typo.Module {
		return &MissingDashes{}
	})
}
