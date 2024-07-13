package outputs

import (
	"fmt"

	"github.com/cybint/athena/interfaces"
)

type Creator func() interfaces.Output

var Outputs = map[string]Creator{}

func Add(name string, creator Creator) {
	Outputs[name] = creator
}

func Get(name string) (Creator, error) {
	if plugin, ok := Outputs[name]; ok {
		return plugin, nil
	}

	return nil, fmt.Errorf("unable to locate outputs/%s plugin", name)
}
