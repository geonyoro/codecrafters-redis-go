package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidNewStreamId(t *testing.T) {
	assert.True(t, IsValidNewStreamId([]int{}, 1, 0))
	assert.False(t, IsValidNewStreamId([]int{1, 0}, 1, 0))
}
