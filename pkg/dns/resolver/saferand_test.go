package resolver

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSafeRandInt31(t *testing.T) {
	tests := []struct {
		name       string
		haveCrypto bool
	}{
		{"Crypto", true},
		{"Time", false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			saferand := newSafeRand(test.haveCrypto)

			n1 := saferand.Int31()
			n2 := saferand.Int31()

			assert.NotEqual(t, n1, n2)
		})
	}
}
