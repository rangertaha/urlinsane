package outputs

import (
	"fmt"

	"github.com/rangertaha/urlinsane"
)

const (
	ENTITY  = "ENTITY"
	DOMAINS = "DOMAINS"
)

type Creator func() urlinsane.Output

var Types = []string{"ENTITY", "DOMAINS"}


var Outputs = map[string]Creator{}

func Add(name string, creator Creator) {
	Outputs[name] = creator
}

func Get(name string) (urlinsane.Output, error) {
	if plugin, ok := Outputs[name]; ok {
		return plugin(), nil
	}

	return nil, fmt.Errorf("unable to locate outputs/%s plugin", name)
}

func All() (mods []urlinsane.Output) {
	for _, plugin := range Outputs {
		mods = append(mods, plugin())
	}
	return
}

func List(IDs ...string) (outputs []urlinsane.Output) {
	for id, output := range Outputs {
		for _, aid := range IDs {
			if id == aid {
				outputs = append(outputs, output())
			}
		}
	}
	for _, aid := range IDs {
		if aid == "all" {
			IDs = []string{}
		}
	}

	if len(IDs) == 0 {
		for _, output := range Outputs {
			outputs = append(outputs, output())
		}
	}

	return
}
