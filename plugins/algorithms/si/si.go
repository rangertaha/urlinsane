package none

import (
	typo "github.com/rangertaha/urlinsane"
	"github.com/rangertaha/urlinsane/plugins/algorithms"
)

type None struct {
	// Code() string
	// Name() string
	// Description() string
	// Fields() []string
	// Exec() func(Result) []Result
}

func (n *None) Code() string {
	return "SI"
}

func (n *None) Name() string {
	return "Subdomain Insertion"
}

func (n *None) Description() string {
	return "inserts common subdomain at the beginning of the domain"
}

func (n *None) Fields() []string {
	return []string{}
}

func (n *None) Exec(typo.Result) (results []typo.Result) {
	return
}

// Register the plugin
func init() {
	algorithms.Add("SI", func() typo.Module {
		return &None{

		}
	})
}
