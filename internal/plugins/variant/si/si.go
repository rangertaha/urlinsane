// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later
// Package si is the subdomain insertion algorithm.
package si

import (
	"github.com/rangertaha/urlinsane/internal/plugins/variant"
)

// subdomainInsertion prepends each known subdomain label to the registrable
// domain. It replaces any existing prefix rather than stacking on top of it,
// which is why it varies the whole key and re-splits itself.
func subdomainInsertion(subdomains []string) variant.Generate {
	return func(name string) []string {
		_, core, suffix := variant.SplitDomain(name)
		if core == "" {
			return nil
		}
		out := make([]string, 0, len(subdomains))
		for _, sub := range subdomains {
			out = append(out, variant.JoinDomain(sub, core, suffix))
		}
		return out
	}
}

// Spec declares the algorithm.
func Spec(subdomains []string) variant.Spec {
	return variant.Spec{
		ID: "si", Title: "Subdomain Insertion", Version: 1,
		Types: []string{variant.TypeDomain},
		Whole: true,
		Gen:   subdomainInsertion(subdomains),
	}
}
