package domains

import (
	"fmt"

	"github.com/rangertaha/urlinsane/internal"
)

// type Infos []internal.Information
// type InformationsOrder struct{ Infos }
// type InfosReverseOrder struct{ Infos }

// func (o Infos) Len() int      { return len(o) }
// func (o Infos) Swap(i, j int) { o[i], o[j] = o[j], o[i] }
// func (l InfosReverseOrder) Less(i, j int) bool {
// 	return l.Infos[i].Order() > l.Infos[j].Order()
// }
// func (l InformationsOrder) Less(i, j int) bool {
// 	return l.Infos[i].Order() < l.Infos[j].Order()
// }

// // sort.Sort(ProcessorOrder{cfgs})


// ----------------
type Creator func() internal.Information

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



// // ----------------
// type Creator func() internal.Information

// var Information = map[string]Creator{}

// func Add(name string, creator Creator) {
// 	Information[name] = creator
// }

// func Get(name string) (Creator, error) {
// 	if plugin, ok := Information[name]; ok {
// 		return plugin, nil
// 	}

// 	return nil, fmt.Errorf("unable to locate information/%s plugin", name)
// }

// func All() (mods []internal.Information) {
// 	for _, plugin := range Information {
// 		mods = append(mods, plugin())
// 	}
// 	return
// }

// func List(IDs ...string) (infos []internal.Information) {
// 	for id, info := range Information {
// 		for _, aid := range IDs {
// 			if id == aid {
// 				infos = append(infos, info())
// 			}
// 		}
// 	}
// 	for _, aid := range IDs {
// 		if aid == "all" {
// 			IDs = []string{}
// 		}
// 	}

// 	if len(IDs) == 0 {
// 		for _, info := range Information {
// 			infos = append(infos, info())
// 		}
// 	}

// 	return
// }
