// Copyright © 2018 CyberSecTech Inc
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

package urlinsane

/*
Package urlinsane implements typosquatting functions

CLI

	urlinsane google.com
    urlinsane google.com -t CO
    urlinsane google.com -t CO -k en2


Example


    package main

	import (
		"fmt"
		"github.com/cybersectech-org/urlinsane"
	)

	func main() {

		conf := urlinsane.BasicConfig{
			Domains:     []string{"google.com"},
			Keyboards:   []string{"en1"},
			Typos:       []string{"co"},
			Funcs:       []string{"ip"},
			Concurrency: 50,
			Format:      "text",
			Verbose:     false,
		}

		urli := urlinsane.New(conf.Config())

		out := urli.Execute()

		for r := range out {
			fmt.Println(r.Live, r.Domain.Domain, r.Typo.Name, r.Data)
		}

	}



### Design

![alt text](design.png "Design Overview")

* Configuration via CLI options.
* Generating typo options to be consumed by workers
* Executing typosquatting algorithm workers concurrently. Defaults to 50 workers at a time.
* Executing extra functions workers on typos to add additional information to the results. Each extra fuc
* Output results in text, json and csv formats.


*/
