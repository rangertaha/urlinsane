package algorithms

import (
	"fmt"

	"github.com/rangertaha/urlinsane"
)

type Creator func() urlinsane.Module

var Algorithms = map[string]Creator{}

func Add(name string, creator Creator) {
	Algorithms[name] = creator
}

func Get(name string) (Creator, error) {
	if plugin, ok := Algorithms[name]; ok {
		return plugin, nil
	}

	return nil, fmt.Errorf("unable to locate outputs/%s plugin", name)
}

func All() (mods []urlinsane.Module) {
	for _, plugin := range Algorithms {
		mods = append(mods, plugin())
	}
	return
}

func List(IDs ...string) (algos []urlinsane.Module) {
	for id, algo := range Algorithms {
		for _, aid := range IDs {
			if id == aid {
				algos = append(algos, algo())
			}
		}
	}
	for _, aid := range IDs {
		if aid == "all" {
			IDs = []string{}
		}
	}

	if len(IDs) == 0 {
		for _, algo := range Algorithms {
			algos = append(algos, algo())
		}
	}

	return
}
