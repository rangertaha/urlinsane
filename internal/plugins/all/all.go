// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

// Package all links in every optional plugin. Importing it for side effects is
// what puts them in a build:
//
//	import _ "github.com/rangertaha/urlinsane/internal/plugins/all"
//
// The shipped operators, algorithms and analyzers are not here. They are
// composed directly by internal/scan through decompose/all, variant/all,
// observe/all and analyze/all, because a registry entry for what always runs
// would be a second place to look for it.
//
// What belongs here is what does *not* always run: a plugin for one registry,
// one forge, one platform. It is empty today, and that is a fact worth stating
// rather than a file worth deleting — without this package a new plugin has no
// obvious place to be linked from, and the first one would be wired into
// cmd/urlinsane by hand, where the second would not find it.
package all
