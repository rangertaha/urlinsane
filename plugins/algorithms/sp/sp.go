package none

import (
	typo "github.com/rangertaha/urlinsane"
	"github.com/rangertaha/urlinsane/plugins/algorithms"
)

type SingularPluralize struct {
	// Code() string
	// Name() string
	// Description() string
	// Fields() []string
	// Exec() func(Result) []Result
}

func (n *SingularPluralize) Code() string {
	return "sp"
}

func (n *SingularPluralize) Name() string {
	return "Singular Pluralize"
}

func (n *SingularPluralize) Description() string {
	return "Singular Pluralise creates a singular domain plural and vice versa"
}

func (n *SingularPluralize) Fields() []string {
	return []string{}
}

func (n *SingularPluralize) Exec(typo.Result) (results []typo.Result) {
	return
}

// Register the plugin
func init() {
	algorithms.Add("sp", func() typo.Module {
		return &SingularPluralize{

		}
	})
}
