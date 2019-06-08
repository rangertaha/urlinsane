// The MIT License (MIT)
//
// Copyright © 2019 Rangertaha <rangertaha@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package typo

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
)

func (urli *URLInsane) outFile() (file *os.File) {
	if urli.file != "" {
		var err error
		file, err = os.OpenFile(urli.file, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		file = os.Stdout
	}
	return
}

func (urli *URLInsane) jsonOutput(in <-chan TypoResult) {

}

func (urli *URLInsane) csvOutput(in <-chan TypoResult) {
	w := csv.NewWriter(urli.outFile())

	live := func(l bool) string {
		if l {
			return "ONLINE"
		} else {
			return " "
		}
	}

	// CSV column headers
	w.Write(urli.headers)

	for v := range in {
		var data []string
		if urli.verbose {
			data = []string{live(v.Live), v.Typo.Name, v.Variant.String(), v.Variant.Suffix}
		} else {
			data = []string{live(v.Live), v.Typo.Code, v.Variant.String(), v.Variant.Suffix}
		}

		// Add a column of data to the results
		for _, head := range urli.headers[4:] {
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

func (urli *URLInsane) stdOutput(in <-chan TypoResult) {
	table := tablewriter.NewWriter(urli.outFile())
	table.SetHeader(urli.headers)
	table.SetBorder(false)

	live := func(l bool) string {
		if l {
			return "ONLINE"
		} else {
			return " "
		}
	}
	for v := range in {
		var data []string
		if urli.verbose {
			data = []string{live(v.Live), v.Typo.Name, v.Variant.String(), v.Variant.Suffix}
		} else {
			data = []string{live(v.Live), v.Typo.Code, v.Variant.String(), v.Variant.Suffix}
		}

		// Add a column of data to the results
		for _, head := range urli.headers[4:] {
			value, ok := v.Data[head]
			if ok {
				data = append(data, value)
			}
		}
		table.Append(data)
	}
	table.Render()
}

func (urli *URLInsane) Output(in <-chan TypoResult) {
	if urli.format == "json" {
		urli.jsonOutput(in)
	}
	if urli.format == "csv" {
		urli.csvOutput(in)
	}
	if urli.format == "text" {
		urli.stdOutput(in)
	}
}
