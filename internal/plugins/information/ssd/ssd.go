// Copyright (C) 2024 Rangertaha
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.
package ssd

// https://github.com/glaslos/ssdeep

import (
	"os"

	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/plugins/information"
	"github.com/rangertaha/urlinsane/pkg/fuzzy/ssdeep"
	log "github.com/sirupsen/logrus"
)

const (
	ORDER       = 9
	CODE        = "ssd"
	NAME        = "SSDeep"
	DESCRIPTION = "SSDeep contents comparison"
)

type Plugin struct {
}

func (n *Plugin) Id() string {
	return CODE
}

func (n *Plugin) Order() int {
	return ORDER
}

func (i *Plugin) Init(conf internal.Config) {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.WarnLevel)
}

func (n *Plugin) Description() string {
	return DESCRIPTION
}

func (n *Plugin) Headers() []string {
	return []string{"SSDEEP"}
}

func (n *Plugin) Exec(in internal.Typo) internal.Typo {
	orig, vari := in.Get()
	// in.Set(orig, vari)

	if vari.Html == "" {
		return in
	}

	hash1, err := ssdeep.FuzzyBytes([]byte(vari.Html))
	if err != nil {
		log.Error(err)
	}
	vari.Ssdeep = hash1

	hash2, _ := ssdeep.FuzzyBytes([]byte(orig.Html))
	if err != nil {
		log.Error(err)
	}

	dist, err := ssdeep.Distance(hash1, hash2)
	if err != nil {
		log.Error(err)
	} else {
		in.SetMeta("SSDEEP", dist)
	}

	in.Set(orig, vari)
	return in
}

// Register the plugin
func init() {
	information.Add(CODE, func() internal.Information {
		return &Plugin{}
	})
}
