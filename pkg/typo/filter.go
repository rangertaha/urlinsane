// The MIT License (MIT)
//
// Copyright © 2018 rangertaha rangertaha@gmail.com
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
	"strings"
)

// FilterREGISTRY The registry for extra functions
var FilterREGISTRY = make(map[string][]Extra)

var onlineFilter = Extra{
	Code:        "LIVE",
	Name:        "Online domians",
	Description: "Show online/live domains only.",
	Exec:        onlineFilterFunc,
	Headers:     []string{"IPv4", "IPv6"},
}

// var onlineFilter = Filter{
// 	Code:        "LIVE",
// 	Name:        "Online domians",
// 	Description: "Show online/live domains only.",
// 	Exec:        onlineFilterFunc,
// 	Requiments:  []string{"IPv4", "IPv6"},
// }

func init() {
	FilterRegister("LIVE", onlineFilter)

	FilterRegister("ALL",
		onlineFilter,
	)
}

// onlineFilterFunc ...
func onlineFilterFunc(tr TypoResult) (results []TypoResult) {
	_, ok := tr.Data["IPv6"]
	if ok {
		if tr.Live {
			results = append(results, TypoResult{Original: tr.Original, Variant: tr.Variant, Typo: tr.Typo, Live: tr.Live, Data: tr.Data})
		}
	}
	return
}

// FilterRegister ...
func FilterRegister(name string, efunc ...Extra) {
	_, registered := FilterREGISTRY[strings.ToUpper(name)]
	if !registered {
		FilterREGISTRY[strings.ToUpper(name)] = efunc
	}
}

// FilterRetrieve ...
func FilterRetrieve(strs ...string) (results []Extra) {
	for _, f := range strs {
		value, ok := FilterREGISTRY[strings.ToUpper(f)]
		if ok {
			results = append(results, value...)
		}
	}
	return
}
