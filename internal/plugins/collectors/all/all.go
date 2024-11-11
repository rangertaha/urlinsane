package all

// register all plugins
import (
	_ "github.com/rangertaha/urlinsane/internal/plugins/collectors/cn"
	_ "github.com/rangertaha/urlinsane/internal/plugins/collectors/geo"
	_ "github.com/rangertaha/urlinsane/internal/plugins/collectors/idn"
	_ "github.com/rangertaha/urlinsane/internal/plugins/collectors/ip"
	_ "github.com/rangertaha/urlinsane/internal/plugins/collectors/mx"
	_ "github.com/rangertaha/urlinsane/internal/plugins/collectors/ns"
	_ "github.com/rangertaha/urlinsane/internal/plugins/collectors/txt"
)
