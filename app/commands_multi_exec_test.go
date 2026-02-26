package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExec_EmptyTransaction(t *testing.T) {
	// should return an array
	ctx := NewTestingContext()
	ctx.ConnState.IsMulti = true
	ret := Exec(ctx, Command{})
	assert.Equal(t, []any{}, ret.EncoderArgs)
}

func TestExec_QueuedCommands(t *testing.T) {
	ctx := NewTestingContext()
	ctx.ConnState.IsMulti = true
	ctx.ConnState.MultiCmds = []MultiCmd{
		{Set, []string{"foo", "xyz"}},
		{Incr, []string{"foo"}},
		{Incr, []string{"bar"}},
		{Get, []string{"bar"}},
	}
	var ret ReturnValue
	ret = Exec(ctx, Command{})

	retEncoderArgs := make([]any, 4)
	for idx, retAnyArg := range ret.EncoderArgs.([]any) {
		ret2Args, ok := retAnyArg.([]byte)
		if !ok {
			panic(fmt.Sprintf("value at outer idx %d is not a ReturnValue", idx))
		}
		retEncoderArgs[idx] = string(ret2Args)
	}
	assert.Equal(t, []any{
		"+OK\r\n",
		"-ERR value is not an integer or out of range\r\n",
		":1\r\n",
		"$1\r\n1\r\n",
	}, retEncoderArgs)
}

func TestExec_WithoutMulti(t *testing.T) {
	ctx := NewTestingContext()
	ret := Exec(ctx, Command{})
	assert.Equal(t, ErrorMultiWithoutExec, ret.EncoderArgs)
}

func TestMulti_ReturnsOK(t *testing.T) {
	ctx := NewTestingContext()
	ret := Multi(ctx, Command{})
	assert.Equal(t, "OK", ret.EncoderArgs)
}

func TestMulti_QueueCommands(t *testing.T) {
	ctx := NewTestingContext()
	ctx.ConnState.IsMulti = true
	var ret ReturnValue
	ret = Multi(ctx, Command{"SET", []string{"foo", "1"}})
	assert.Equal(t, "QUEUED", ret.EncoderArgs)

	ret = Multi(ctx, Command{"INCR", []string{"foo"}})
	assert.Equal(t, "QUEUED", ret.EncoderArgs)

	// verify that the foo value doesn't exist becauase it's not been set yet
	ret = Get(ctx, Command{"GET", []string{"foo"}})
	assert.Equal(t, nil, ret.EncoderArgs)
}
