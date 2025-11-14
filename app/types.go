package main

import (
	"io"
	"sync"
)

type Variable struct {
	Value              string
	SetAt              int64
	ExpiryMilliseconds int64
}

type ListVariable struct {
	Values []string
	Mu     sync.Mutex
}

type (
	SequenceKV map[string]string // maps a key to a value
	MillisVal  struct {
		Map  map[string]SequenceKV // maps a sequenceId to Entry
		Keys []string              // stored as sorted list of ints
		Last int                   // latest
	}
	Stream struct {
		Map  map[string]*MillisVal // maps streamId to the StreamValue
		Keys []string              // stored as sorted list of ints
		Last int                   // latest
	}
)

func NewMillisVal() *MillisVal {
	return &MillisVal{
		Map:  make(map[string]SequenceKV),
		Last: -1,
	}
}

func NewStream() *Stream {
	return &Stream{
		Map:  make(map[string]*MillisVal),
		Last: -1,
	}
}

type Context struct {
	Conn  io.ReadWriteCloser
	State *State
}

type State struct {
	VariableMap map[string]Variable
	ListMap     map[string]*ListVariable
	StreamMap   map[string]*Stream
}

func NewState() *State {
	vMap := make(map[string]Variable)
	lMap := make(map[string]*ListVariable)
	sMap := make(map[string]*Stream)
	return &State{
		VariableMap: vMap,
		ListMap:     lMap,
		StreamMap:   sMap,
	}
}

type XRangeReturn struct {
	ID string
	KV map[string]string
}
