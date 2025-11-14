package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseInput(t *testing.T) {
	output, err := ParseInput("*2\r\n$4\r\nECHO\r\n$3\r\nhey\r\n")
	assert.Nil(t, err)
	assert.Equal(t, output, Command{
		Command: "ECHO",
		Args: []string{
			"hey",
		},
	})
}

func TestParseXReadArgs(t *testing.T) {
	actual := ParseXReadArgs([]string{
		"streams", "stream_key", "other_stream_key", "0-0", "0-1",
	})
	assert.Equal(t, map[string]string{
		"stream_key":       "0-0",
		"other_stream_key": "0-1",
	}, actual)
}
