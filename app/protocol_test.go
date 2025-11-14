package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRArrayAllStrings(t *testing.T) {
	assert.Equal(t, RArray([]any{}), []byte("*0\r\n"))
	assert.Equal(t, RArray([]any{"hello", "world"}), []byte("*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n"))
}

func TestRArrayInt(t *testing.T) {
	assert.Equal(t, RArray([]any{1, 2}), []byte("*2\r\n:1\r\n:2\r\n"))
}

func TestRArrayMixedIntString(t *testing.T) {
	assert.Equal(t, RArray([]any{1, "hello"}), []byte("*2\r\n:1\r\n$5\r\nhello\r\n"))
}

func TestRArrayMixedIntStringArray(t *testing.T) {
	input := RArray([]any{
		1,
		"foo",
		[]any{
			"hello",
			"world",
		},
	})
	assert.Equal(t, input, []byte("*3\r\n:1\r\n$3\r\nfoo\r\n*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n"))
}

func TestRArrayTwoNested(t *testing.T) {
	actual := RArray([]any{
		[]any{
			1,
			2,
			3,
		},
		[]any{
			"hello",
			"world",
		},
	})
	expected := []byte("*2\r\n*3\r\n:1\r\n:2\r\n:3\r\n*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n")
	fmt.Printf("actual: %#v\n", string(actual))
	fmt.Printf("expected: %#v\n", string(expected))
	assert.Equal(t, expected, actual)
}
