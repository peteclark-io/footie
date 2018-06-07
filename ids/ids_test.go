package ids

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewIDs(t *testing.T) {
	s1 := NewID()
	s2 := NewID()
	assert.NotEqual(t, s1, s2)
}
