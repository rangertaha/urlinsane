package pi

import (
	"github.com/rangertaha/urlinsane"
	"github.com/rangertaha/urlinsane/plugins/algorithms"
)

const CODE = "pi"

type PeriodInsertion struct {
	types []string
}

func (n *PeriodInsertion) Code() string {
	return CODE
}
func (n *PeriodInsertion) IsType(str string) bool {
	return algorithms.IsType(n.types, str)
}

func (n *PeriodInsertion) Name() string {
	return "PeriodInsertion"
}

func (n *PeriodInsertion) Description() string {
	return "Inserting periods in the target name"
}

func (n *PeriodInsertion) Fields() []string {
	return []string{}
}

func (n *PeriodInsertion) Headers() []string {
	return []string{}
}

func (n *PeriodInsertion) Exec(in urlinsane.Typo) (out []urlinsane.Typo) {
	out = append(out, in)
	return
}

// Register the plugin
func init() {
	algorithms.Add(CODE, func() urlinsane.Algorithm {
		return &PeriodInsertion{
			types: []string{algorithms.ENTITY, algorithms.DOMAIN},
		}
	})
}
