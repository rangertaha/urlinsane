package databases

import (
	"fmt"

	"github.com/rangertaha/urlinsane/internal"
)

type Creator func() internal.Database

var Databases = map[string]Creator{}

func Add(name string, creator Creator) {
	Databases[name] = creator
}

func Get(name string) (internal.Database, error) {
	if plugin, ok := Databases[name]; ok {
		return plugin(), nil
	}

	return nil, fmt.Errorf("unable to locate databases/%s plugin", name)
}

func All() (mods []internal.Database) {
	for _, plugin := range Databases {
		mods = append(mods, plugin())
	}
	return
}

func List(IDs ...string) (databases []internal.Database) {
	for id, output := range Databases {
		for _, aid := range IDs {
			if id == aid {
				databases = append(databases, output())
			}
		}
	}
	for _, aid := range IDs {
		if aid == "all" {
			IDs = []string{}
		}
	}

	if len(IDs) == 0 {
		for _, database := range Databases {
			databases = append(databases, database())
		}
	}

	return
}
