package information

import (
	"fmt"

	"github.com/rangertaha/urlinsane"
)

const (
	ENTITY  = "ENTITY"
	DOMAINS = "DOMAINS"
)

type Creator func() urlinsane.Information

var Types = []string{"ENTITY", "DOMAINS"}

var Information = map[string]Creator{}

func Add(name string, creator Creator) {
	Information[name] = creator
}

func Get(name string) (Creator, error) {
	if plugin, ok := Information[name]; ok {
		return plugin, nil
	}

	return nil, fmt.Errorf("unable to locate information/%s plugin", name)
}

func All() (mods []urlinsane.Information) {
	for _, plugin := range Information {
		mods = append(mods, plugin())
	}
	return
}

func List(IDs ...string) (infos []urlinsane.Information) {
	for id, info := range Information {
		for _, aid := range IDs {
			if id == aid {
				infos = append(infos, info())
			}
		}
	}
	for _, aid := range IDs {
		if aid == "all" {
			IDs = []string{}
		}
	}

	if len(IDs) == 0 {
		for _, info := range Information {
			infos = append(infos, info())
		}
	}

	return
}
