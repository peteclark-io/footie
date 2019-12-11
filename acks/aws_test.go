package acks

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReceive(t *testing.T) {
	to := `\"in+group6789match6789pedro6789@shoreditch.football\" <in+group6789match6789pedro6789@shoreditch.football>`
	matches := addressRegex.FindStringSubmatch(to)

	assert.Len(t, matches, 4)
	assert.Equal(t, matches[1], "group6789")
	assert.Equal(t, matches[2], "match6789")
	assert.Equal(t, matches[3], "pedro6789")
}
