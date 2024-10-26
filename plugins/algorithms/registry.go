package algorithms

import (
	"fmt"

	"github.com/rangertaha/urlinsane"
)

const (
	ENTITY = "ENTITY"
	DOMAINS = "DOMAINS"
)


type Creator func() urlinsane.Algorithm

var Types = []string{"ENTITY", "DOMAINS"}

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

func All() (mods []urlinsane.Algorithm) {
	for _, plugin := range Algorithms {
		mods = append(mods, plugin())
	}
	return
}

func List(IDs ...string) (algos []urlinsane.Algorithm) {
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


func IsType(types []string, other string) bool {
	for _, typ := range types {
		if typ == other {
			return true
		}
	}
	return false
}