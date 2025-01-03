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
	"slices"
	"time"

	"github.com/rangertaha/urlinsane/internal/dataset"
)

func Domains(language, file string) (err error) {
	// fmt.Println("Importing domains...")
	fmt.Printf("Importing domains from %s\n", file)
	var domains []*dataset.Domain
	for _, line := range Extract(file) {
		for _, domain := range line {
			var dm dataset.Domain
			dataset.Data.FirstOrInit(&dm, dataset.Domain{Name: domain})
			domains = append(domains, &dm)
		}
	}

	for chunk := range slices.Chunk(domains, 1000) {
		fmt.Printf("Saving %v domains \n", len(chunk))
		dataset.Data.Save(&chunk)
		time.Sleep(time.Second / 4)
	}

	return
}


func Suffixes(language, file string) (err error) {
	fmt.Printf("Importing suffixes from %s\n", file)
	var suffixes []*dataset.Suffix
	for _, line := range Extract(file) {
		for _, suffix := range line {
			var s dataset.Suffix
			dataset.Data.FirstOrInit(&s, dataset.Suffix{Name: suffix})
			suffixes = append(suffixes, &s)
		}
	}

	for chunk := range slices.Chunk(suffixes, 1000) {
		fmt.Printf("Saving %v suffixes \n", len(chunk))
		dataset.Data.Save(&chunk)
		time.Sleep(time.Second / 4)
	}

	return
}


func Prefixes(language, file string) (err error) {
	fmt.Printf("Importing prefixes from %s\n", file)
	var prefixes []*dataset.Prefix
	for _, line := range Extract(file) {
		for _, domain := range line {
			var p dataset.Prefix
			dataset.Data.FirstOrInit(&p, dataset.Prefix{Name: domain})
			prefixes = append(prefixes, &p)
		}
	}

	for chunk := range slices.Chunk(prefixes, 1000) {
		fmt.Printf("Saving %v prefixes \n", len(chunk))
		dataset.Data.Save(&chunk)
		time.Sleep(time.Second / 4)
	}

	return
}
