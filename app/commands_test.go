package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type DummyConn struct {
	Data []byte
}

func (d *DummyConn) Read(p []byte) (n int, err error) {
	copySize := copy(d.Data, p)
	return copySize, nil
}

func (d *DummyConn) Write(p []byte) (n int, err error) {
	d.Data = append(d.Data, p...)
	return len(d.Data), nil
}

func (d *DummyConn) Close() error {
	return nil
}

func TestLrange_NegativeNos(t *testing.T) {
	dConn := DummyConn{
		Data: []byte{},
	}
	ctx := &Context{
		Conn:  &dConn,
		State: NewState(),
	}
	lmap := *ctx.State.ListMap
	lvar := ListVariable{
		Values: []string{"a", "b", "c", "d", "e"},
	}
	lmap["list"] = &lvar
	testCases := []struct {
		description string
		index1      string
		index2      string
		output      []string
	}{
		{
			description: "last 2 elements",
			index1:      "-2",
			index2:      "-1",
			output:      []string{"d", "e"},
		},
		{
			description: "all elements except the last 2",
			index1:      "0",
			index2:      "-3",
			output:      []string{"a", "b", "c"},
		},
		{
			description: "all elements with negative indexes",
			index1:      "-5",
			index2:      "-1",
			output:      []string{"a", "b", "c", "d", "e"},
		},
		{
			description: "all elements with out of bound negatives",
			index1:      "-6",
			index2:      "-1",
			output:      []string{"a", "b", "c", "d", "e"},
		},
		{
			description: "all elements with out of bound negatives and bigger",
			index1:      "-6",
			index2:      "-7",
			output:      []string{},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			cmd := Command{
				Command: "LRANGE",
				Args:    []string{"list", testCase.index1, testCase.index2},
			}

			retVal := Lrange(ctx, cmd)
			assert.Equal(t, testCase.output, retVal.EncoderArgs)
		})
	}
}

func TestLrangeBasic(t *testing.T) {
	dConn := DummyConn{
		Data: []byte{},
	}
	ctx := &Context{
		Conn:  &dConn,
		State: NewState(),
	}
	lmap := *ctx.State.ListMap
	lvar := ListVariable{
		Values: []string{"a", "b", "c", "d", "e"},
	}
	lmap["list"] = &lvar
	cmd := Command{
		Command: "LRANGE",
		Args:    []string{"list", "1", "2"},
	}

	output := Lrange(ctx, cmd)
	assert.Equal(t, []string{"b", "c"}, output.EncoderArgs)
}

func TestConvertLrangeIndex(t *testing.T) {
	for _, testCase := range []struct {
		listSize int
		index    int
		output   uint
	}{
		{listSize: 5, index: 0, output: 0},
		{listSize: 5, index: 4, output: 4},
		{listSize: 5, index: 5, output: 4},
		{listSize: 5, index: -1, output: 4},
		{listSize: 5, index: -5, output: 0},
		{listSize: 5, index: -6, output: 0},
	} {
		t.Run(fmt.Sprintf("%#v", testCase), func(t *testing.T) {
			output := convertLrangeIndex(testCase.listSize, testCase.index)
			assert.Equal(t, testCase.output, output)
		})
	}
}
