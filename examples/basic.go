package main

import (
	"fmt"

	"github.com/cybint/urlinsane/pkg/typo"
)

func main() {

	conf := typo.BasicConfig{
		Domains:     []string{"google.com"},
		Keyboards:   []string{"en1"},
		Typos:       []string{"co"},
		Funcs:       []string{"ip"},
		Concurrency: 50,
		Format:      "text",
		Verbose:     false,
	}

	urli := typo.New(conf.Config())

	out := urli.Stream()

	for r := range out {
		fmt.Println(r.Variant.Live, r.Variant.Domain, r.Typo.Name, r.Data)
	}
}
