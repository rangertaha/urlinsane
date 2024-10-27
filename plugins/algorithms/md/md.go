package md

import (
	"github.com/rangertaha/urlinsane"
	"github.com/rangertaha/urlinsane/plugins/algorithms"
)

const CODE = "md"

type MissingDot struct {
	types []string
}

func (n *MissingDot) Id() string {
	return CODE
}
func (n *MissingDot) IsType(str string) bool {
	return algorithms.IsType(n.types, str)
}

func (n *MissingDot) Name() string {
	return "Missing Dot"
}

func (n *MissingDot) Description() string {
	return "Created by omitting a dot from the name"
}

func (n *MissingDot) Fields() []string {
	return []string{}
}

func (n *MissingDot) Headers() []string {
	return []string{}
}

func (n *MissingDot) Exec(typo urlinsane.Typo) (typos []urlinsane.Typo) {
	// fmt.Println("Pre:", typo)
	for _, variant := range missingCharFunc(typo.Original().Repr(), ".") {
		// fmt.Println("Variant:", variant)

		if typo.Original().Repr() != variant {
			// newTypo := typo.Copy()
			// newTypo := typo.NewVariant(variant)
			// fmt.Println(variant)
			typos = append(typos, typo.New(variant))
		}
	}
	// fmt.Println("Post:", typo)

	// typos = append(typos, typo)

	return
}

// Register the plugin
func init() {
	algorithms.Add(CODE, func() urlinsane.Algorithm {
		return &MissingDot{
			types: []string{algorithms.DOMAIN},
		}
	})
}

// func missingDotFunc(tc Result) (results []Result) {
// 	for _, str := range missingCharFunc(tc.Original.String(), ".") {
// 		if tc.Original.Domain != str {
// 			dm := Domain{tc.Original.Subdomain, str, tc.Original.Suffix, Meta{}, false}
// 			results = append(results, Result{Original: tc.Original, Variant: dm, Typo: tc.Typo, Data: tc.Data})
// 		}
// 	}
// 	dm := Domain{tc.Original.Subdomain, strings.Replace(tc.Original.Domain, ".", "", -1), tc.Original.Suffix, Meta{}, false}
// 	results = append(results, Result{Original: tc.Original, Variant: dm, Typo: tc.Typo, Data: tc.Data})
// 	return results
// }

// missingCharFunc removes a character one at a time from the string.
// For example, wwwgoogle.com and www.googlecom
func missingCharFunc(str, character string) (results []string) {
	for i, char := range str {
		if character == string(char) {
			results = append(results, str[:i]+str[i+1:])
		}
	}
	return
}
