package main

import (
	"fmt"
	"strings"
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
	lmap := ctx.State.ListMap
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
	lmap := ctx.State.ListMap
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

func TestLPush(t *testing.T) {
	ctx := &Context{
		Conn:  &DummyConn{},
		State: NewState(),
	}
	lMap := ctx.State.ListMap
	list := []string{"1", "2", "3"}
	lMap["list"] = &ListVariable{
		Values: list,
	}
	cmd := Command{
		"LPUSH",
		[]string{"list", "a", "b", "c"},
	}
	ret := Lpush(ctx, cmd)
	assert.Equal(t, 6, ret.EncoderArgs)
}

func TestLlen(t *testing.T) {
	ctx := &Context{
		Conn:  &DummyConn{},
		State: NewState(),
	}
	lMap := ctx.State.ListMap
	list := []string{"1", "2", "3"}
	lMap["list"] = &ListVariable{
		Values: list,
	}
	cmd := Command{
		"LLEN",
		[]string{"list"},
	}
	ret := Llen(ctx, cmd)
	assert.Equal(t, 3, ret.EncoderArgs)
}

func TestXadd(t *testing.T) {
	ctx := &Context{
		Conn:  &DummyConn{},
		State: NewState(),
	}
	millis := "1526919030474"
	sequence := "0"
	streamId := fmt.Sprintf("%s-%s", millis, sequence)
	cmd := Command{
		"XADD",
		[]string{"stream_key", streamId, "temperature", "36", "humidity", "95"},
	}
	ret := Xadd(ctx, cmd)
	assert.Equal(t, streamId, ret.EncoderArgs)
	sMap := ctx.State.StreamMap
	stream, ok := sMap["stream_key"]
	assert.True(t, ok) // ensure the stream key exists
	millisVal := stream.Map[millis]
	entryMap := millisVal.Map[sequence]
	assert.Equal(t, "95", entryMap["humidity"])
	assert.Equal(t, "36", entryMap["temperature"])
}

func TestXaddWithPartial(t *testing.T) {
	ctx := &Context{
		Conn:  &DummyConn{},
		State: NewState(),
	}
	millis := "1526919030474"
	sequence := "*"
	streamId := fmt.Sprintf("%s-%s", millis, sequence)
	cmd := Command{
		"XADD",
		[]string{"stream_key", streamId, "temperature", "36", "humidity", "95"},
	}
	ret := Xadd(ctx, cmd)
	assert.True(t, strings.HasPrefix(ret.EncoderArgs.(string), millis))
	sMap := ctx.State.StreamMap
	stream, ok := sMap["stream_key"]
	assert.True(t, ok) // ensure the stream key exists
	millisVal := stream.Map[millis]
	entryMap := millisVal.Map["0"]
	assert.Equal(t, "36", entryMap["temperature"])
	assert.Equal(t, "95", entryMap["humidity"])

	// adding another makes it 1
	cmd = Command{
		"XADD",
		[]string{"stream_key", streamId, "temperature", "50", "humidity", "100"},
	}
	ret = Xadd(ctx, cmd)
	// these haven't changed
	entryMap = millisVal.Map["0"]
	assert.Equal(t, "36", entryMap["temperature"])
	assert.Equal(t, "95", entryMap["humidity"])
	// for the new entry
	entryMap = millisVal.Map["1"]
	assert.Equal(t, "50", entryMap["temperature"])
	assert.Equal(t, "100", entryMap["humidity"])
}
