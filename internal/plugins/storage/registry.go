package storage

import (
	"fmt"

	"github.com/rangertaha/urlinsane/internal"
)

type Creator func() internal.Storage

var Stores = map[string]Creator{}

func Add(name string, creator Creator) {
	Stores[name] = creator
}

func Get(name string) (internal.Storage, error) {
	if plugin, ok := Stores[name]; ok {
		return plugin(), nil
	}

	return nil, fmt.Errorf("unable to locate storage/%s plugin", name)
}

func All() (mods []internal.Storage) {
	for _, plugin := range Stores {
		mods = append(mods, plugin())
	}
	return
}

func List(IDs ...string) (stores []internal.Storage) {
	for id, output := range Stores {
		for _, aid := range IDs {
			if id == aid {
				stores = append(stores, output())
			}
		}
	}
	for _, aid := range IDs {
		if aid == "all" {
			IDs = []string{}
		}
	}

	if len(IDs) == 0 {
		for _, store := range Stores {
			stores = append(stores, store())
		}
	}

	return
}
