package bf

// // bitsquattingFunc relies on random bit- errors to redirect connections
// // intended for popular domains
// func bitsquattingFunc(tc Result) (results []Result) {
// 	// TOOO: need to improve.
// 	masks := []int{1, 2, 4, 8, 16, 32, 64, 128}
// 	charset := make(map[string][]string)
// 	for _, board := range tc.Keyboards {
// 		for _, alpha := range board.Language.Graphemes {
// 			for _, mask := range masks {
// 				new := int([]rune(alpha)[0]) ^ mask
// 				for _, a := range board.Language.Graphemes {
// 					if string(a) == string(new) {
// 						charset[string(alpha)] = append(charset[string(alpha)], string(new))
// 					}
// 				}
// 			}
// 		}
// 	}

// 	for d, dchar := range tc.Original.Domain {
// 		for _, char := range charset[string(dchar)] {

// 			dnew := tc.Original.Domain[:d] + string(char) + tc.Original.Domain[d+1:]
// 			dm := Domain{tc.Original.Subdomain, dnew, tc.Original.Suffix, Meta{}, false}
// 			results = append(results, Result{Original: tc.Original, Variant: dm, Typo: tc.Typo, Data: tc.Data})
// 		}
// 	}
// 	return
// }

import (
	"github.com/rangertaha/urlinsane"
	"github.com/rangertaha/urlinsane/plugins/algorithms"
)

const CODE = "bf"

type BitFlipping struct {
	types []string
}

func (n *BitFlipping) Id() string {
	return CODE
}
func (n *BitFlipping) IsType(str string) bool {
	return algorithms.IsType(n.types, str)
}

func (n *BitFlipping) Name() string {
	return "Bit Flipping"
}

func (n *BitFlipping) Description() string {
	return "Relies on random bit-errors to redirect connections"
}

func (n *BitFlipping) Fields() []string {
	return []string{}
}

func (n *BitFlipping) Headers() []string {
	return []string{}
}

func (n *BitFlipping) Exec(in urlinsane.Typo) (out []urlinsane.Typo) {
	out = append(out, in)
	return
}

// Register the plugin
func init() {
	algorithms.Add(CODE, func() urlinsane.Algorithm {
		return &BitFlipping{
			types: []string{algorithms.ENTITY, algorithms.DOMAIN},
		}
	})
}
