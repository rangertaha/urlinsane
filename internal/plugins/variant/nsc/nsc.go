// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later
// Package nsc is the namespace confusion algorithm.
package nsc

import (
	"github.com/rangertaha/urlinsane/internal/plugins/variant"
	"strings"
)

// namespaceConfusion moves a name between the namespacing conventions the
// registries use — npm's "@org/pkg", a repo's "org/pkg", and the flat
// "org-pkg" — which is how a scoped package gets impersonated by an unscoped
// one and vice versa.
func namespaceConfusion(name string) []string {
	switch {
	case strings.HasPrefix(name, "@") && strings.Contains(name, "/"):
		org, pkg, _ := strings.Cut(strings.TrimPrefix(name, "@"), "/")
		return []string{pkg, org + "-" + pkg, org + "/" + pkg, org + pkg}
	case strings.Contains(name, "/"):
		org, pkg, _ := strings.Cut(name, "/")
		return []string{pkg, org + "-" + pkg, "@" + org + "/" + pkg}
	case strings.Contains(name, "-"):
		org, pkg, _ := strings.Cut(name, "-")
		return []string{"@" + org + "/" + pkg, org + "/" + pkg}
	}
	return nil
}

// Spec declares the algorithm.
func Spec() variant.Spec {
	return variant.Spec{
		ID: "nsc", Title: "Namespace Confusion", Version: 1,
		Types: []string{variant.TypePackage, variant.TypeRepo},
		Whole: true,
		Gen:   namespaceConfusion,
	}
}
