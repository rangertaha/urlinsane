// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

// Package tos is the token order algorithm.
//
// A name made of more than one word can be written with those words in another
// order and still read as the same thing: shop-online and online-shop,
// 2024example and example2024, cloud-backup and backup-cloud. Nothing is
// misspelled, so the character-level generators cannot reach it — `cs`
// transposes two adjacent *characters*, which turns shop-online into
// shpo-online, not into online-shop.
//
// It matters most where multi-word names are the convention rather than the
// exception, which is every package registry: node-fetch and fetch-node,
// python-dateutil and dateutil-python. A developer who half-remembers a name
// remembers the words, not their order.
//
// ail-typo-squatting calls this ChangeOrder.
package tos

import (
	"github.com/rangertaha/urlinsane/internal/plugins/variant"
	"github.com/rangertaha/urlinsane/pkg/typo"
)

// Spec declares the algorithm.
//
// No type restriction: a hyphenated domain, a scoped package and a handle all
// carry word order, and the tokenizer finds the boundaries in each without
// being told which it is looking at.
func Spec() variant.Spec {
	return variant.Spec{
		ID: "tos", Title: "Token Order Swap", Version: 1,
		Gen: func(name string) []string { return typo.TokenOrderSwap(name) },
	}
}
