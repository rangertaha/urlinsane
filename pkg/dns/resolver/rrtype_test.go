package resolver

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRRTypeString(t *testing.T) {
	testTable := []struct {
		name string
		have RRtype
		want string
	}{
		{name: "A", have: TypeA, want: "A"},
		{name: "AAAA", have: TypeAAAA, want: "AAAA"},
		{name: "CNAME", have: TypeCNAME, want: "CNAME"},
		{name: "TXT", have: TypeTXT, want: "TXT"},
		{name: "MX", have: TypeMX, want: "MX"},
		{name: "NS", have: TypeNS, want: "NS"},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			got := test.have.String()
			assert.Equal(t, test.want, got)
		})
	}
}
