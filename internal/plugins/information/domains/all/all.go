package all

// register all plugins
import (
	_ "github.com/rangertaha/urlinsane/internal/plugins/information/domains/cn"
	_ "github.com/rangertaha/urlinsane/internal/plugins/information/domains/geo"
	_ "github.com/rangertaha/urlinsane/internal/plugins/information/domains/har"
	_ "github.com/rangertaha/urlinsane/internal/plugins/information/domains/idn"
	_ "github.com/rangertaha/urlinsane/internal/plugins/information/domains/ip"
	_ "github.com/rangertaha/urlinsane/internal/plugins/information/domains/mx"
	_ "github.com/rangertaha/urlinsane/internal/plugins/information/domains/ns"
	_ "github.com/rangertaha/urlinsane/internal/plugins/information/domains/ssd"
	_ "github.com/rangertaha/urlinsane/internal/plugins/information/domains/txt"
	// _ "github.com/rangertaha/urlinsane/internal/plugins/information/domains/web"
	_ "github.com/rangertaha/urlinsane/internal/plugins/information/domains/wi"
	_ "github.com/rangertaha/urlinsane/internal/plugins/information/domains/bn"
)
