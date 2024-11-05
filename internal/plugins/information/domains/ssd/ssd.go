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
	log "github.com/sirupsen/logrus"
)

const (
	CODE        = "ssd"
	NAME        = "SSDeep"
	DESCRIPTION = "SSDeep contents comparison"
)

type None struct {
}

func (n *None) Id() string {
	return CODE
}

func (i *None) Init(conf internal.Config) {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.WarnLevel)
}

func (n *None) Description() string {
	return DESCRIPTION
}

func (n *None) Headers() []string {
	return []string{NAME}
}

func (n *None) Exec(in internal.Typo) internal.Typo {
	// originPage := in.Original().Get("HTML").(string)

	// hash1, err := ssdeep.FuzzyBytes([]byte(Data1))
	// if err != nil {
	// 	log.Error(err)
	// }
	// hash2, _ := ssdeep.FuzzyBytes([]byte(Data2))
	// if err != nil {
	// 	log.Error(err)
	// }
	// in.Variant().Add("HASH1", hash1)
	// in.Variant().Add("HASH2", hash2)

	// dist, err := ssdeep.Distance(hash1, hash2)
	// if err != nil {
	// 	log.Error(err)
	// } else {
	// 	in.Variant().Add(NAME, dist)
	// }

	return in
}

// Register the plugin
func init() {
	information.Add(CODE, func() internal.Information {
		return &None{}
	})
}
