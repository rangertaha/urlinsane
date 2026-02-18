package outputs

import (
	"fmt"
	"time"

	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/db"
)

type Creator func() internal.Output

var Outputs = map[string]Creator{}

func Add(name string, creator Creator) {
	Outputs[name] = creator
}

func Get(name string) (internal.Output, error) {
	if plugin, ok := Outputs[name]; ok {
		return plugin(), nil
	}

	return nil, fmt.Errorf("unable to locate outputs/%s plugin", name)
}

func All() (mods []internal.Output) {
	for _, plugin := range Outputs {
		mods = append(mods, plugin())
	}
	return
}

func List(IDs ...string) (outputs []internal.Output) {
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

type Plugin struct {
	Summary string
	ID      string
	Config  internal.Config
	Domains []*db.Domain
	Elapsed time.Duration
	Started time.Time
	Offline int64
	Online  int64
	Total   int64
}

func (p *Plugin) Id() string {
	return p.ID
}

func (p *Plugin) Description() string {
	return p.Summary
}

func (p *Plugin) Init(conf internal.Config) {
	p.Started = time.Now()
	p.Config = conf
}

func (p *Plugin) Report() {
	p.Elapsed = time.Since(p.Started)
	summary := map[string]string{
		"  TIME:":  p.Elapsed.String(),
		"  TOTAL:": fmt.Sprintf("%d", p.Total),
	}
	if len(p.Config.Collectors()) > 0 {
		summary[text.FgGreen.Sprintf("%s", "  LIVE:")] = fmt.Sprintf("%d", p.Online)
		summary[text.FgRed.Sprintf("%s", "  OFFLINE")] = fmt.Sprintf("%d", p.Offline)
	}

	fmt.Println("")
	for k, v := range summary {
		fmt.Printf("%s %s   ", k, v)
	}
	fmt.Println("")
}
