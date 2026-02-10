package main

import (
	"io"
	"strconv"
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
		Keys []string              // stored as sorted list
		Last int                   // latest
	}
	Stream struct {
		Map  map[string]*MillisVal // maps streamId to the StreamValue
		Keys []string              // stored as sorted list
		Last int                   // latest
	}
)

func (s *Stream) GetLastMillisSequence() (int, int) {
	toMillis := s.Last
	millisVal := s.Map[strconv.Itoa(toMillis)]
	toSequence := millisVal.Last
	return toMillis, toSequence
}

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

type ConnState struct {
	IsMulti   bool
	MultiCmds []MultiCmd
}

func NewConnState() *ConnState {
	return &ConnState{
		MultiCmds: make([]MultiCmd, 0),
	}
}

type Context struct {
	Conn      io.ReadWriteCloser
	ConnState *ConnState
	State     *State
}

type MultiCmd struct {
	Callable func(ctx *Context, cmd Command) ReturnValue
	Args     []string
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

func (x XRangeReturn) ToRArray() (r []any) {
	r = make([]any, 0)
	r = append(r, x.ID)
	kvArray := make([]any, 0)
	for key, val := range x.KV {
		kvArray = append(kvArray, key)
		kvArray = append(kvArray, val)
	}
	r = append(r, kvArray)
	return r
}

type XReadReturn struct {
	Stream  string
	Entries []XRangeReturn
}

func (x XReadReturn) ToRArray() (r []any) {
	r = make([]any, 0)
	r = append(r, x.Stream)
	kvArray := make([]any, 0)
	for _, val := range x.Entries {
		kvArray = append(kvArray, val.ToRArray())
	}
	r = append(r, kvArray)
	return r
}
