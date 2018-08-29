package urlinsane

import (
	"encoding/csv"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"os"
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
			data = []string{live(v.Live), v.Typo.Name, v.Domain.String(), v.Domain.Suffix}
		} else {
			data = []string{live(v.Live), v.Typo.Code, v.Domain.String(), v.Domain.Suffix}
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
			data = []string{live(v.Live), v.Typo.Name, v.Domain.String(), v.Domain.Suffix}
		} else {
			data = []string{live(v.Live), v.Typo.Code, v.Domain.String(), v.Domain.Suffix}
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
