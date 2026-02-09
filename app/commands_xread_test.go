package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseXReadArgs_OnlyStreams(t *testing.T) {
	args, err := ParseXReadArgs([]string{
		"streams", "stream_key", "other_stream_key", "0-0", "0-1",
	})
	assert.Nil(t, err)
	assert.Equal(t, map[string]string{
		"stream_key":       "0-0",
		"other_stream_key": "0-1",
	}, args.Streams)
}

func TestParseXReadArgs_WCount(t *testing.T) {
	args, err := ParseXReadArgs([]string{
		"count", "2", "streams", "stream_key", "other_stream_key", "0-0", "0-1",
	})
	assert.Nil(t, err)
	assert.Equal(t, map[string]string{
		"stream_key":       "0-0",
		"other_stream_key": "0-1",
	}, args.Streams)
}

func TestParseXReadArgs_WBlock(t *testing.T) {
	args, err := ParseXReadArgs([]string{
		"block", "100", "streams", "stream_key", "other_stream_key", "0-0", "0-1",
	})
	assert.Nil(t, err)
	assert.Equal(t, map[string]string{
		"stream_key":       "0-0",
		"other_stream_key": "0-1",
	}, args.Streams)
}

func TestParseXReadArgs_WBlock_WCount(t *testing.T) {
	args, err := ParseXReadArgs([]string{
		"count", "2", "block", "100", "streams", "stream_key", "other_stream_key", "0-0", "0-1",
	})
	assert.Nil(t, err)
	assert.Equal(t, map[string]string{
		"stream_key":       "0-0",
		"other_stream_key": "0-1",
	}, args.Streams)
}

func TestParseXReadArgs_WBlock_WCount_Error(t *testing.T) {
	_, err := ParseXReadArgs([]string{
		"count", "x", "block", "100", "streams", "stream_key", "other_stream_key", "0-0", "0-1",
	})
	assert.NotNil(t, err)
}
