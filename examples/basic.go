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

	out := urli.Stream()

	for r := range out {
		fmt.Println(r.Live, r.Variant.Domain, r.Typo.Name, r.Data)
	}
}
