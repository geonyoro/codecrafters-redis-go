package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidNewStreamId(t *testing.T) {
	s := NewStream()
	isValid, err := s.IsNewStreamIdValid("1", "0")
	assert.True(t, isValid)
	assert.Nil(t, err)

	// add these keys to the stream
	s.AddIdWithKV("1", "0", map[string]string{"foo": "bar"})

	// rerun tests
	isValid, err = s.IsNewStreamIdValid("1", "0")
	assert.Nil(t, err)
	assert.False(t, isValid)
}
