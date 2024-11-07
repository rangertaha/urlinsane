package information

import (
	"fmt"

	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/plugins/information/domains"
	"github.com/rangertaha/urlinsane/internal/plugins/information/packages"
	"github.com/rangertaha/urlinsane/internal/plugins/information/usernames"
)
<<<<<<< HEAD

// const (
// 	ENTITY  = "ENTITY"
// 	DOMAINS = "DOMAINS"
// )

type Creator func() internal.Information

// var Types = []string{"ENTITY", "DOMAINS"}

=======

type Creator func() internal.Information

>>>>>>> develop
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

func All() (mods []internal.Information) {
	for _, plugin := range Information {
		mods = append(mods, plugin())
	}
	return
}

func List(IDs ...string) (infos []internal.Information) {
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

func ListType(ttype int, IDs ...string) (infos []internal.Information) {
	if internal.DOMAIN == ttype {
		return domains.List(IDs...)
	}
	if internal.PACKAGE == ttype {
		return packages.List(IDs...)
	}
	if internal.NAME == ttype {
		return usernames.List(IDs...)
	}
	return
}
