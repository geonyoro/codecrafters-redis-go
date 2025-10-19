package main

import (
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

// func TestLrange_NegativeNos(t *testing.T) {
// 	dConn := DummyConn{
// 		Data: []byte{},
// 	}
// 	ctx := &Context{
// 		Conn:  &dConn,
// 		State: NewState(),
// 	}
// 	lmap := *ctx.State.ListMap
// 	lvar := ListVariable{
// 		Values: []string{"a", "b", "c", "d", "e"},
// 	}
// 	lmap["list"] = &lvar
// 	cmd := Command{
// 		Command: "LRANGE",
// 		Args:    []string{"list", "-2", "-1"},
// 	}
//
// 	Lrange(ctx, cmd)
// 	assert.Equal(t, []byte("*2\r\n$1\r\nd\r\n$1\r\ne\r\n"), dConn.Data)
// }

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

	Lrange(ctx, cmd)
	assert.Equal(t, []byte("*2\r\n$1\r\nb\r\n$1\r\nc\r\n"), dConn.Data)
}
