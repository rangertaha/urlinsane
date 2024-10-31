package all

// register all plugins
import (
	_ "github.com/rangertaha/urlinsane/internal/plugins/information/domains"
	_ "github.com/rangertaha/urlinsane/internal/plugins/information/emails"
	_ "github.com/rangertaha/urlinsane/internal/plugins/information/packages"
	_ "github.com/rangertaha/urlinsane/internal/plugins/information/usernames"
)
