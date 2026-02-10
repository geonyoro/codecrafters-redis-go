package main

import (
	"fmt"
	"strings"
	"testing"
	"time"

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

func TestIncr_WithExistantValue(t *testing.T) {
	ctx := &Context{
		Conn:  &DummyConn{},
		State: NewState(),
	}
	key := "x"
	(ctx.State.VariableMap)[key] = Variable{
		Value: "1",
		SetAt: time.Now().UnixMilli(),
	}

	// perform the increase
	cmd := Command{
		"INCR",
		[]string{key},
	}
	ret := Incr(ctx, cmd)
	assert.Equal(t, 2, ret.EncoderArgs)

	// again yields 3
	ret = Incr(ctx, cmd)
	assert.Equal(t, 3, ret.EncoderArgs)
}

func TestIncr_WithInexistantValue(t *testing.T) {
	ctx := &Context{
		Conn:  &DummyConn{},
		State: NewState(),
	}
	key := "x"
	// perform the increase
	cmd := Command{
		"INCR",
		[]string{key},
	}
	ret := Incr(ctx, cmd)
	assert.Equal(t, 1, ret.EncoderArgs)
}

func TestIncr_WithNonNumericValue(t *testing.T) {
	ctx := &Context{
		Conn:  &DummyConn{},
		State: NewState(),
	}
	key := "x"
	(ctx.State.VariableMap)[key] = Variable{
		Value: "y",
		SetAt: time.Now().UnixMilli(),
	}

	// perform the increase
	cmd := Command{
		"INCR",
		[]string{key},
	}
	ret := Incr(ctx, cmd)
	assert.Equal(t, ErrorIncr, ret.EncoderArgs)
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
