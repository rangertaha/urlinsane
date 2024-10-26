package pi

import (
	"github.com/rangertaha/urlinsane"
	"github.com/rangertaha/urlinsane/plugins/algorithms"
)

type PeriodInsertion struct {
	types []string
}

func (n *PeriodInsertion) Code() string {
	return "pi"
}
func (n *PeriodInsertion) IsType(str string) bool {
	return algorithms.IsType(n.types, str)
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

func (n *PeriodInsertion) Exec(urlinsane.Typo) (results []urlinsane.Typo) {
	return
}

// Register the plugin
func init() {
	algorithms.Add("pi", func() urlinsane.Algorithm {
		return &PeriodInsertion{
			types: []string{algorithms.ENTITY, algorithms.DOMAINS},
		}
	})
}
