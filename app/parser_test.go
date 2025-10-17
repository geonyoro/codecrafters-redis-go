package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseInput(t *testing.T) {
	output := ParseInput("*2\r\n$4\r\nECHO\r\n$3\r\nhey\r\n")
	assert.Equal(t, output, Command{
		Command: "ECHO",
		Args: []string{
			"hey",
		},
	})
}
