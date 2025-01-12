package all

// register all plugins
import (
	_ "github.com/rangertaha/urlinsane/internal/plugins/collectors/cn"
	_ "github.com/rangertaha/urlinsane/internal/plugins/collectors/geo"
	_ "github.com/rangertaha/urlinsane/internal/plugins/collectors/idn"

	// _ "github.com/rangertaha/urlinsane/internal/plugins/collectors/img"
	// _ "github.com/rangertaha/urlinsane/internal/plugins/collectors/bn"
	_ "github.com/rangertaha/urlinsane/internal/plugins/collectors/ip"
	_ "github.com/rangertaha/urlinsane/internal/plugins/collectors/mx"
	_ "github.com/rangertaha/urlinsane/internal/plugins/collectors/ns"
	_ "github.com/rangertaha/urlinsane/internal/plugins/collectors/ptr"
	_ "github.com/rangertaha/urlinsane/internal/plugins/collectors/txt"
	// _ "github.com/rangertaha/urlinsane/internal/plugins/collectors/web"
	_ "github.com/rangertaha/urlinsane/internal/plugins/collectors/wi"
)
