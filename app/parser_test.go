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

func TestParseCliArgs_Listener(t *testing.T) {
	args := ParseCliArgs([]string{})
	assert.Equal(t, "127.0.0.1", args.Host)
	assert.Equal(t, 6379, args.Port)
}

func TestParseCliArgs_ReplicaOf(t *testing.T) {
	args := ParseCliArgs([]string{"--replicaof", "localhost 6379"})
	assert.Equal(t, "localhost 6379", args.ReplicaOf)
}
