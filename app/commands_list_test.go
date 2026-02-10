package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLrange_NegativeNos(t *testing.T) {
	ctx := NewTestingContext()
	lmap := ctx.State.ListMap
	lvar := ListVariable{
		Values: []string{"a", "b", "c", "d", "e"},
	}
	lmap["list"] = &lvar
	testCases := []struct {
		description string
		index1      string
		index2      string
		output      []any
	}{
		{
			description: "last 2 elements",
			index1:      "-2",
			index2:      "-1",
			output:      []any{"d", "e"},
		},
		{
			description: "all elements except the last 2",
			index1:      "0",
			index2:      "-3",
			output:      []any{"a", "b", "c"},
		},
		{
			description: "all elements with negative indexes",
			index1:      "-5",
			index2:      "-1",
			output:      []any{"a", "b", "c", "d", "e"},
		},
		{
			description: "all elements with out of bound negatives",
			index1:      "-6",
			index2:      "-1",
			output:      []any{"a", "b", "c", "d", "e"},
		},
		{
			description: "all elements with out of bound negatives and bigger",
			index1:      "-6",
			index2:      "-7",
			output:      []any{},
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
	ctx := NewTestingContext()
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
	assert.Equal(t, []any{"b", "c"}, output.EncoderArgs)
}

func TestLPush(t *testing.T) {
	ctx := NewTestingContext()
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
	ctx := NewTestingContext()
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
