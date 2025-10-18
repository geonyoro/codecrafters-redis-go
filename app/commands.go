package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func ExecuteCommand(ctx *Context, cmd Command) bool {
	cmdFunc, ok := CmdFuncMap[strings.ToUpper(cmd.Command)]
	if ok {
		cmdFunc(ctx, cmd)
		return true
	}
	fmt.Println("Failed to find cmd for", cmd.Command)
	return false
}

func Echo(ctx *Context, cmd Command) {
	output := strings.Join(cmd.Args, " ")
	ctx.Conn.Write(RBulkString(output))
}

func Ping(ctx *Context, cmd Command) {
	ctx.Conn.Write(RSimpleString("PONG"))
}

func Set(ctx *Context, cmd Command) {
	key := cmd.Args[0]
	value := cmd.Args[1]

	expiryMilliseconds := int64(-1)
	for i := range len(cmd.Args) / 2 {
		if i == 0 {
			continue
		}
		idx := i * 2
		arg := cmd.Args[idx]
		if arg == "EX" || arg == "PX" {
			intArgString := cmd.Args[idx+1]
			mult := 1
			if arg == "EX" {
				mult = 1000
			}
			if intArg, err := strconv.Atoi(intArgString); err == nil {
				expiryMilliseconds = int64(mult * intArg)
			}
		}
	}

	(*ctx.State.VariableMap)[key] = Variable{
		Value:              value,
		SetAt:              time.Now().UnixMilli(),
		ExpiryMilliseconds: expiryMilliseconds,
	}
	ctx.Conn.Write(RSimpleString("OK"))
}

func Get(ctx *Context, cmd Command) {
	key := cmd.Args[0]
	value, ok := (*ctx.State.VariableMap)[key]
	if ok {
		isExpired := false
		nowMillis := time.Now().UnixMilli()
		if value.ExpiryMilliseconds > 0 {
			expiresAt := value.SetAt + value.ExpiryMilliseconds
			if expiresAt <= nowMillis {
				isExpired = true
			}
		}
		if !isExpired {
			ctx.Conn.Write(RBulkString(value.Value))
			return
		}
	}
	ctx.Conn.Write(RNullBulkString())
}

func Rpush(ctx *Context, cmd Command) {
	listName := cmd.Args[0]
	listMap := *ctx.State.ListMap
	list, ok := listMap[listName]
	if !ok {
		list = ListVariable{}
	}
	for i := 1; i < len(cmd.Args); i++ {
		newValue := cmd.Args[i]
		list.Values = append(list.Values, newValue)
	}
	ctx.Conn.Write(RInteger(len(list.Values)))
}

var CmdFuncMap = map[string]func(ctx *Context, cmd Command){
	"ECHO":  Echo,
	"PING":  Ping,
	"SET":   Set,
	"GET":   Get,
	"RPUSH": Rpush,
}
