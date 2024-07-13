package algorithms

import (
	"fmt"

	typo "github.com/rangertaha/urlinsane"
)

type Creator func() typo.Module

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

func All() (mods []typo.Module) {
	for _, plugin := range Algorithms {
		mods = append(mods, plugin())
	}
	return
}