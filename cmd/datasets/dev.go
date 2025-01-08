// Copyright 2024 Rangertaha. All Rights Reserved.
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

package main

import (
	"fmt"

	"github.com/rangertaha/urlinsane/internal/config"
	"github.com/rangertaha/urlinsane/internal/dataset"
	"github.com/urfave/cli/v2"
)

var DevCmd = cli.Command{
	Name:                   "dev",
	Usage:                  "Dev",
	UsageText:              "",
	UseShortOptionHandling: true,
	Flags:                  importFlags,
	Action: func(cCtx *cli.Context) error {
		return Dev(cCtx)
	},
}

func Dev(cli *cli.Context) error {
	// var tlds dataset.Suffix
	config.New(config.CliOptions(cli))

	tlds := []dataset.Suffix{}
	dataset.DB.Find(&tlds)

	for _, tld := range tlds {
		fmt.Println(tld.Name)
	}

	// result.RowsAffected // returns found records count, equals `len(users)`
	// result.Error        // returns error

	return nil
}
