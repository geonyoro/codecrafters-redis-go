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

