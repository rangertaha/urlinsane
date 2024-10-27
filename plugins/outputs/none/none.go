package none

import (
	"fmt"

	"github.com/rangertaha/urlinsane"
	"github.com/rangertaha/urlinsane/config"
	"github.com/rangertaha/urlinsane/plugins/outputs"
)

const CODE = "none"

type None struct {
	file string
}

func (n *None) Init(conf config.Config) {
	n.file = conf.File
}

func (n *None) Id() string {
	return CODE
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
