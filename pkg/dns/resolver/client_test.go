package resolver

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	resolvers := []string{}

	client := New(resolvers, 3, 10, 5)

	assert.NotNil(t, client)
}
