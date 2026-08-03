package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/rangertaha/urlinsane/pkg/kb"
)

func main() {
	for _, id := range os.Args[1:] {
		l, err := kb.Get(id)
		if err != nil {
			fmt.Println(id, "ERR", err)
			continue
		}
		fmt.Println("#### ", id, l.Name, l.Languages())
		for _, row := range l.Rows() {
			var b []string
			for _, k := range row {
				t := k.Base()
				if t == "" {
					continue
				}
				b = append(b, t)
			}
			fmt.Println("  ROW:", strings.Join(b, " "))
		}
		seen := map[string]bool{}
		for _, row := range l.Rows() {
			for _, k := range row {
				t := k.Base()
				if t == "" || seen[t] {
					continue
				}
				seen[t] = true
				fmt.Printf("  ADJ %s : %s\n", t, strings.Join(l.Adjacent(t), " "))
			}
		}
	}
}
