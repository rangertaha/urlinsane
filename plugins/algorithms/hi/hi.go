package hi

import (
	"github.com/rangertaha/urlinsane"
	"github.com/rangertaha/urlinsane/plugins/algorithms"
)

const CODE = "hi"

type DashInsertion struct {
	types []string
}

func (n *DashInsertion) Code() string {
	return CODE
}
func (n *DashInsertion) IsType(str string) bool {
	return algorithms.IsType(n.types, str)
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
	algorithms.Add(CODE, func() urlinsane.Algorithm {
		return &DashInsertion{
			types: []string{algorithms.ENTITY, algorithms.DOMAINS},
		}
	})
}
