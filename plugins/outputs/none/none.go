package none

import (
	"fmt"

	"github.com/rangertaha/urlinsane"
	"github.com/rangertaha/urlinsane/plugins/outputs"
)

const CODE = "none"

type None struct {
	rtype string
	file  string
}

func (n *None) Code() string {
	return CODE
}

func (n *None) Set(typ, filepath string) {
	n.rtype = typ
	n.file = filepath
}

func (n *None) Description() string {
	return ""
}

func (n *None) Write(in urlinsane.Typo) {
	fmt.Println(in)
}

// Register the plugin
func init() {
	outputs.Add(CODE, func() urlinsane.Output {
		return &None{}
	})
}
