package text

import (
	"fmt"

	"github.com/rangertaha/urlinsane"
	"github.com/rangertaha/urlinsane/plugins/outputs"
)

const CODE = "text"

type Text struct {
	rtype string
	file  string
}

func (n *Text) Code() string {
	return CODE
}

func (n *Text) Set(typ string, filepath string) {
	n.rtype = typ
	n.file = filepath

}

func (n *Text) Description() string {
	return "Text outputs one record per line and is the default"
}

func (n *Text) Write(in urlinsane.Typo) {
	fmt.Println(in)
}

// Register the plugin
func init() {
	outputs.Add(CODE, func() urlinsane.Output {
		return &Text{}
	})
}
