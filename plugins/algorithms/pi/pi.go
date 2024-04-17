package none

import (
	typo "github.com/rangertaha/urlinsane"
	"github.com/rangertaha/urlinsane/plugins/algorithms"
)

type PeriodInsertion struct {
	// Code() string
	// Name() string
	// Description() string
	// Fields() []string
	// Exec() func(Result) []Result
}

func (n *PeriodInsertion) Code() string {
	return "pi"
}

func (n *PeriodInsertion) Name() string {
	return "PeriodInsertion"
}

func (n *PeriodInsertion) Description() string {
	return "Inserting periods in the target domain"
}

func (n *PeriodInsertion) Fields() []string {
	return []string{}
}

func (n *PeriodInsertion) Headers() []string {
	return []string{}
}

func (n *PeriodInsertion) Exec(typo.Result) (results []typo.Result) {
	return
}

// Register the plugin
func init() {
	algorithms.Add("pi", func() typo.Module {
		return &PeriodInsertion{}
	})
}
