package analyzers

import (
	"fmt"

	"github.com/rangertaha/urlinsane/internal"
)

type Creator func() internal.Analyzer

var Analyzers = map[string]Creator{}

func Add(name string, creator Creator) {
	Analyzers[name] = creator
}

func Get(name string) (internal.Analyzer, error) {
	if plugin, ok := Analyzers[name]; ok {
		return plugin(), nil
	}

	return nil, fmt.Errorf("unable to locate outputs/%s plugin", name)
}

func All() (mods []internal.Analyzer) {
	for _, plugin := range Analyzers {
		mods = append(mods, plugin())
	}
	return
}

func List(IDs ...string) (analyzers []internal.Analyzer) {
	for id, analyzer := range Analyzers {
		for _, aid := range IDs {
			if id == aid {
				analyzers = append(analyzers, analyzer())
			}
		}
	}
	for _, aid := range IDs {
		if aid == "all" {
			IDs = []string{}
		}
	}

	if len(IDs) == 0 {
		for _, analyzer := range Analyzers {
			analyzers = append(analyzers, analyzer())
		}
	}

	return
}
