package none

import (
	"github.com/rangertaha/urlinsane"
	"github.com/rangertaha/urlinsane/plugins/algorithms"
)

type DashInsertion struct {
	// Code() string
	// Name() string
	// Description() string
	// Fields() []string
	// Exec() func(Result) []Result
}

func (n *DashInsertion) Code() string {
	return "hi"
}

func (n *DashInsertion) Name() string {
	return "Dash Insertion"
}

func (n *DashInsertion) Description() string {
	return "Inserting hyphens in the target domain"
}

func (n *DashInsertion) Fields() []string {
	return []string{}
}

func (n *DashInsertion) Headers() []string {
	return []string{}
}

func (n *DashInsertion) Exec(urlinsane.Typo) (results []urlinsane.Typo) {
	return
}

// Register the plugin
func init() {
	algorithms.Add("hi", func() urlinsane.Algorithm {
		return &DashInsertion{}
	})
}
