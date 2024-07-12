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

// Filters ...
var Filters = NewRegistry()

var onlineFilter = Module{
	Code:        "LIVE",
	Name:        "Online domians",
	Description: "Show online/live domains only.",
	Exe:         onlineFilterFunc,
	Fields:      []string{"IPv4", "IPv6"},
}

func init() {
	Filters.Set("LIVE", onlineFilter)

	Filters.Set("ALL",
		onlineFilter,
	)
}

// onlineFilterFunc ...
func onlineFilterFunc(tr Result) (results []Result) {
	_, ok := tr.Data["IPv6"]
	if ok {
		if tr.Variant.Live {
			results = append(results, Result{Original: tr.Original, Variant: tr.Variant, Typo: tr.Typo, Data: tr.Data})
		}
	}
	return
}
