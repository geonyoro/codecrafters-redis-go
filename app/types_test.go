package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidNewStreamId(t *testing.T) {
	s := NewStream()
	assert.True(t, s.IsNewStreamIdValid(1, 0))
	// add these keys to the stream
	s.AddIdWithKV("1", "0", map[string]string{"foo": "bar"})
	assert.False(t, s.IsNewStreamIdValid(1, 0))
}
