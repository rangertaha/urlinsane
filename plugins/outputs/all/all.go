package all

import (
	// register plugin
	_ "github.com/rangertaha/urlinsane/plugins/outputs/csv" 
	_ "github.com/rangertaha/urlinsane/plugins/outputs/tsv" 
	_ "github.com/rangertaha/urlinsane/plugins/outputs/html" 
	_ "github.com/rangertaha/urlinsane/plugins/outputs/md" 
	_ "github.com/rangertaha/urlinsane/plugins/outputs/text" 
)