package md

// missingDotFunc typos are created by omitting a dot from the domain. For example, wwwgoogle.com and www.googlecom

import (
	"github.com/rangertaha/urlinsane"
	"github.com/rangertaha/urlinsane/plugins/algorithms"
	"github.com/rangertaha/urlinsane/utils/nlp"
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
	for _, variant := range nlp.MissingCharFunc(typo.Original().Repr(), ".") {
		if typo.Original().Repr() != variant {
			typos = append(typos, typo.New(variant))
		}
	}
	return
}

// Register the plugin
func init() {
	algorithms.Add(CODE, func() urlinsane.Algorithm {
		return &MissingDot{
			types: []string{algorithms.ENTITY, algorithms.DOMAIN},
		}
	})
}
