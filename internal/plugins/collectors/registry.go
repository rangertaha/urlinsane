package collectors

import (
	"fmt"

	"github.com/rangertaha/urlinsane/internal"
	log "github.com/sirupsen/logrus"
)

var HEAVY = []string{"img", "bn"}

type Creator func() internal.Collector

var Collector = map[string]Creator{}

func Add(name string, creator Creator) {
	Collector[name] = creator
}

func Get(name string) (Creator, error) {
	if plugin, ok := Collector[name]; ok {
		return plugin, nil
	}

	return nil, fmt.Errorf("unable to locate collector/%s plugin", name)
}

func All() (mods []internal.Collector) {
	for _, plugin := range Collector {
		mods = append(mods, plugin())
	}
	return
}

func List(IDs ...string) (infos []internal.Collector) {
	for id, info := range Collector {
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
		for _, info := range Collector {
			infos = append(infos, info())
		}
	}

	return
}

type Plugin struct {
	Conf      internal.Config
	DependsOn []string
	Summary   string
	Title     string
	Code      string
	Num       int
	Log       *log.Entry
}

func (p *Plugin) Id() string {
	return p.Code
}

func (p *Plugin) Init(c internal.Config) {
	p.Log = log.WithFields(log.Fields{"plugin": p.Code, "method": "Exec"})
	p.Conf = c
}

func (p *Plugin) Order() int {
	return p.Num
}

func (p *Plugin) Name() string {
	return p.Title
}

func (p *Plugin) Dependencies() []string {
	return p.DependsOn
}

func (p *Plugin) Description() string {
	return p.Summary
}
