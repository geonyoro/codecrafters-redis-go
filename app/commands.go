package main

import (
	"fmt"
	"strings"
)

func ExecuteCommand(ctx *Context, cmd Command) bool {
	cmdFunc, ok := CmdFuncMap[cmd.Command]
	if ok {
		cmdFunc(ctx, cmd)
		return true
	}
	fmt.Println("Failed to find cmd for", cmd.Command)
	return false
}

func Echo(ctx *Context, cmd Command) {
	output := strings.Join(cmd.Args, " ")
	ctx.Conn.Write(BulkString(output))
}

func Ping(ctx *Context, cmd Command) {
	ctx.Conn.Write(SimpleString("PONG"))
}

func Set(ctx *Context, cmd Command) {
	key := cmd.Args[0]
	value := cmd.Args[1]
	(*ctx.VariableMap)[key] = Variable{Value: value, Expiry: -1}
	ctx.Conn.Write(SimpleString("OK"))
}

func Get(ctx *Context, cmd Command) {
	key := cmd.Args[0]
	value, ok := (*ctx.VariableMap)[key]
	if ok {
		ctx.Conn.Write(BulkString(value.Value))
		return
	}
	ctx.Conn.Write(BulkString("OK"))
}

var CmdFuncMap = map[string]func(ctx *Context, cmd Command){
	"ECHO": Echo,
	"PING": Ping,
	"SET":  Set,
	"GET":  Get,
}
