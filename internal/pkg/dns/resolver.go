// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

// Package dns holds the resolver the observation operators query through.
//
// It is deliberately one variable now. The package also carried a domain
// parser and a set of name-permutation helpers; the parser read the public
// suffix list out of the dataset database, which made an operator's output
// depend on whether a database happened to be open, so variant.SplitDomain
// reads the compiled-in list instead and the parser lost its callers.
package dns

import "net"

// Resolver is the DNS resolver the observation operators use.
//
// It once had a SetResolver to point it at custom servers, wired from a
// --nameservers flag. The flag is gone, so the setter had no caller; querying
// several resolvers is still worth doing and is on the roadmap, but it will
// want a resolver per operator call rather than one package-level swap.
var Resolver = net.DefaultResolver
