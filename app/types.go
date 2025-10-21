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

type Context struct {
	Conn  io.ReadWriteCloser
	State *State
}

type State struct {
	VariableMap *map[string]Variable
	ListMap     *map[string]*ListVariable
}

func NewState() *State {
	vMap := make(map[string]Variable)
	lMap := make(map[string]*ListVariable)
	return &State{
		VariableMap: &vMap,
		ListMap:     &lMap,
	}
}
