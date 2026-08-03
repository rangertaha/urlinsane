//go:build ignore

package main

import (
	"fmt"
	"os"

	"github.com/rangertaha/urlinsane/pkg/kb"
)

func main() {
	for _, tag := range []string{"kk", "ky", "tt", "ru"} {
		fmt.Printf("### ByLanguage(%s):\n", tag)
		for _, e := range kb.ByLanguage(tag) {
			fmt.Printf("  %s %s\n", e.ID, e.Name)
		}
	}
	for _, id := range os.Args[1:] {
		l := kb.MustGet(id)
		fmt.Printf("=== %s (%s)\n", l.ID, l.Name)
		seen := map[string]bool{}
		for _, k := range l.Keys {
			for _, t := range k.Texts() {
				if len([]rune(t)) != 1 || seen[t] {
					continue
				}
				seen[t] = true
				fmt.Printf("%s\t%v\n", t, l.Adjacent(t))
			}
		}
	}
}
