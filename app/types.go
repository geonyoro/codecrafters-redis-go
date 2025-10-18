package main

import "net"

type Variable struct {
	Value              string
	SetAt              int64
	ExpiryMilliseconds int64
}

type Context struct {
	Conn        net.Conn
	VariableMap *map[string]Variable
}

func NewContext(conn net.Conn, variableMap *map[string]Variable) *Context {
	return &Context{
		conn,
		variableMap,
	}
}
