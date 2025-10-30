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
	Entry  map[string]string // maps a key to a value
	Stream struct {
		Entries   map[string]Entry // maps stream Id to the ValueSet
		LastEntry []int
	}
)

func IsValidNewStreamId(lastEntry []int, millis, sequence int) bool {
	if len(lastEntry) == 0 {
		return true
	}
	if lastEntry[0] > millis {
		return false
	}
	if lastEntry[1] >= sequence {
		return false
	}
	return true
}

type Context struct {
	Conn  io.ReadWriteCloser
	State *State
}

type State struct {
	VariableMap *map[string]Variable
	ListMap     *map[string]*ListVariable
	StreamMap   *map[string]*Stream
}

func NewState() *State {
	vMap := make(map[string]Variable)
	lMap := make(map[string]*ListVariable)
	sMap := make(map[string]*Stream)
	return &State{
		VariableMap: &vMap,
		ListMap:     &lMap,
		StreamMap:   &sMap,
	}
}
