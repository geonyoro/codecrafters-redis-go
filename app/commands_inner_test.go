package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestXrangeInnerWithSequence(t *testing.T) {
	type testCase struct {
		From string
		To   string
		Ret  []XRangeReturn
	}
	d := DummyConn{}
	ctx := &Context{
		Conn:  &d,
		State: NewState(),
	}
	stream := NewStream()
	streamId := "myStream"
	ctx.State.StreamMap[streamId] = stream
	// assumes this method has been thoroughly tested
	stream.AddIdWithKV("1", "0", map[string]string{
		"a": "1",
		"b": "2",
	})
	stream.AddIdWithKV("1", "1", map[string]string{
		"x": "9",
		"y": "0",
	})
	stream.AddIdWithKV("1", "2", map[string]string{
		"a": "3",
		"b": "4",
	})
	stream.AddIdWithKV("2", "3", map[string]string{
		"c": "5",
		"d": "6",
	})
	stream.AddIdWithKV("4", "0", map[string]string{
		"a": "7",
		"b": "8",
	})
	stream.AddIdWithKV("4", "4", map[string]string{
		"a": "2",
	})
	for _, testCase := range []testCase{
		{
			From: "-", To: "1-1", Ret: []XRangeReturn{
				{
					ID: "1-0", KV: map[string]string{"a": "1", "b": "2"},
				},
				{
					ID: "1-1", KV: map[string]string{"x": "9", "y": "0"},
				},
			},
		},
		{
			From: "1-0", To: "1-1", Ret: []XRangeReturn{
				{
					ID: "1-0", KV: map[string]string{"a": "1", "b": "2"},
				},
				{
					ID: "1-1", KV: map[string]string{"x": "9", "y": "0"},
				},
			},
		},
		{
			From: "1-1", To: "1-2", Ret: []XRangeReturn{
				{
					ID: "1-1", KV: map[string]string{"x": "9", "y": "0"},
				},
				{
					ID: "1-2", KV: map[string]string{"a": "3", "b": "4"},
				},
			},
		},
		{
			From: "2-0", To: "3-0", Ret: []XRangeReturn{
				{
					ID: "2-3", KV: map[string]string{"c": "5", "d": "6"},
				},
			},
		},
		{
			From: "2-0", To: "4-0", Ret: []XRangeReturn{
				{
					ID: "2-3", KV: map[string]string{"c": "5", "d": "6"},
				},
				{
					ID: "4-0", KV: map[string]string{"a": "7", "b": "8"},
				},
			},
		},
		{
			From: "3-0", To: "+", Ret: []XRangeReturn{
				{
					ID: "4-0", KV: map[string]string{"a": "7", "b": "8"},
				},
				{
					ID: "4-4", KV: map[string]string{"a": "2"},
				},
			},
		},
	} {
		description := fmt.Sprintf("From:%s To:%s", testCase.From, testCase.To)
		t.Run(description, func(t *testing.T) {
			ret := xRangeInner(ctx, streamId, testCase.From, testCase.To)
			assert.Equal(t, testCase.Ret, ret)
		})
	}
}
