package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExec_EmptyTransaction(t *testing.T) {
	// should return an array
	ctx := &Context{
		State: NewState(),
	}
	ctx.State.IsMulti = true
	ret := Exec(ctx, Command{})
	assert.Equal(t, []any{}, ret.EncoderArgs)
}

func TestExec_WithoutMulti(t *testing.T) {
	dConn := DummyConn{
		Data: []byte{},
	}
	ctx := &Context{
		Conn:  &dConn,
		State: NewState(),
	}
	ret := Exec(ctx, Command{})
	assert.Equal(t, ErrorMultiWithoutExec, ret.EncoderArgs)
}

func TestMulti_ReturnsOK(t *testing.T) {
	dConn := DummyConn{
		Data: []byte{},
	}
	ctx := &Context{
		Conn:  &dConn,
		State: NewState(),
	}
	ret := Multi(ctx, Command{})
	assert.Equal(t, "OK", ret.EncoderArgs)
}

func TestMulti_QueuesCommands(t *testing.T) {
	dConn := DummyConn{
		Data: []byte{},
	}
	ctx := &Context{
		Conn:  &dConn,
		State: NewState(),
	}
	ctx.State.IsMulti = true
	var ret ReturnValue
	ret = Multi(ctx, Command{"SET", []string{"foo", "1"}})
	assert.Equal(t, "QUEUED", ret.EncoderArgs)

	ret = Multi(ctx, Command{"INCR", []string{"foo"}})
	assert.Equal(t, "QUEUED", ret.EncoderArgs)
}
