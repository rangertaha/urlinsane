package languages

import "github.com/rangertaha/urlinsane/internal"

type Language func() internal.Language

type Keyboard func() internal.Keyboard

// var Languages = map[string]Creator{}
var KEYBOARDS = map[string]Keyboard{}
var LANGUAGES = map[string]Language{}

func AddLanguage(name string, lang Language) {
	LANGUAGES[name] = lang
}

func AddKeyboard(name string, kboard Keyboard) {
	KEYBOARDS[name] = kboard
}

// func Get(name string) (Creator, error) {
// 	if plugin, ok := Languages[name]; ok {
// 		return plugin, nil
// 	}

// 	return nil, fmt.Errorf("unable to locate outputs/%s plugin", name)
// }

// func Languages() (mods []internal.Language) {
// 	for _, plugin := range LANGUAGES {
// 		mods = append(mods, plugin())
// 	}
// 	return
// }

func Languages(IDs ...string) (infos []internal.Language) {
	for id, info := range LANGUAGES {
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
		for _, info := range LANGUAGES {
			infos = append(infos, info())
		}
	}

	return
}

func Keyboards(IDs ...string) (infos []internal.Keyboard) {
	for id, info := range KEYBOARDS {
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
		for _, info := range KEYBOARDS {
			infos = append(infos, info())
		}
	}

	return
}
