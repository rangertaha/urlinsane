// Copyright (C) 2024  Rangertaha <rangertaha@gmail.com>
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
package typo

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/olekukonko/tablewriter"
)

func (urli *Typosquatting) outFile() (file *os.File) {
	if urli.config.file != "" {
		var err error
		file, err = os.OpenFile(urli.config.file, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		file = os.Stdout
	}
	return
}

func (urli *Typosquatting) jsonOutput(in <-chan Result) {
	for r := range in {
		if urli.config.verbose {
			json, err := json.MarshalIndent(r, "", "  ")
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(string(json))
		} else {
			json, err := json.MarshalIndent(r, "", "  ")
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(string(json))
		}

	}
}

func (urli *Typosquatting) csvOutput(in <-chan Result) {
	w := csv.NewWriter(urli.outFile())

	live := func(l bool) string {
		if l {
			return "ONLINE"
		}
		return " "
	}

	// CSV column headers
	w.Write(urli.config.headers)

	for v := range in {
		var data []string
		if urli.config.verbose {
			data = []string{live(v.Variant.Live), v.Typo.Name, v.Variant.String(), v.Variant.Suffix}
		} else {
			data = []string{live(v.Variant.Live), v.Typo.Code, v.Variant.String(), v.Variant.Suffix}
		}

		// Add a column of data to the results
		for _, head := range urli.config.headers[4:] {
			value, ok := v.Data[head]
			if ok {
				data = append(data, value)
			}
		}
		if err := w.Write(data); err != nil {
			fmt.Println("Error writing record to csv:", err)
		}
	}
	w.Flush()

	if err := w.Error(); err != nil {
		fmt.Println(err)
	}
}

func (urli *Typosquatting) stdOutput(in <-chan Result) {
	table := tablewriter.NewWriter(urli.outFile())
	table.SetHeader(urli.config.headers)
	table.SetBorder(false)

	live := func(l bool) string {
		if l {
			return "\033[32mONLINE"
		}
		return "\033[39m"
	}
	for v := range in {
		var data []string
		if urli.config.verbose {
			data = []string{live(v.Variant.Live), v.Typo.Name, v.Variant.String(), v.Variant.Suffix}
		} else {
			data = []string{live(v.Variant.Live), v.Typo.Code, v.Variant.String(), v.Variant.Suffix}
		}

		// Add a column of data to the results
		for _, head := range urli.config.headers[4:] {
			value, ok := v.Data[head]
			if ok {
				data = append(data, value)
			}
		}
		table.Append(data)
	}
	table.Render()
	fmt.Println("\033[39m")
}
