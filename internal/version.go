// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later
package internal

const (
	// VERSION format is loosely based on
	// [Semantic Versioning](https://semver.org/spec/v2.0.0.html).
	VERSION = "0.9.0"

	DEBUG = false
)

// COMMIT is the git revision, set at link time by the Makefile:
//
//	-ldflags "-X github.com/rangertaha/urlinsane/internal.COMMIT=$(git rev-parse --short HEAD)"
//
// A var rather than a const because only a var can be set by -X, and "unknown"
// rather than empty so a `go build` with no ldflags still says something true.
var COMMIT = "unknown"

const (

	// LOGO made as ASCII graphics
	BANNER = `
 _   _  ____   _      ___
| | | ||  _ \ | |    |_ _| _ __   ___   __ _  _ __    ___
| | | || |_) || |     | | | '_ \ / __| / _' || '_ \  / _ \
| |_| ||  _ < | |___  | | | | | |\__ \| (_| || | | ||  __/
 \___/ |_| \_\|_____||___||_| |_||___/ \__,_||_| |_| \___|   
 v%s
    
ENTITY:     %s
LANGUAGES:  %s
KEYBOARDS:  %s
ALGORITHMS: %s
COLLECTORS: %s
OUTPUT:     %s
TIME:       %s

`
)
