package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type ReturnValue struct {
	Encoder     func(arg any) []byte
	EncoderArgs any
}

func ExecuteCommand(ctx *Context, cmd Command) bool {
	var returnVal ReturnValue
	isCmdMulti := strings.ToUpper(cmd.Command) == "MULTI"
	isCmdExec := strings.ToUpper(cmd.Command) == "EXEC"

	if ctx.ConnState.IsMulti || isCmdMulti {
		if isCmdExec {
			returnVal = Exec(ctx, cmd)
			fmt.Printf("ExecuteCommand: %#v\n", returnVal)
		} else {
			returnVal = Multi(ctx, cmd)
		}
	} else {
		cmdFunc, ok := CmdFuncMap[strings.ToUpper(cmd.Command)]
		if !ok {
			fmt.Println("Failed to find cmd for", cmd.Command)
			return false
		}
		returnVal = cmdFunc(ctx, cmd)
	}
	encodedVal := returnVal.Encoder(returnVal.EncoderArgs)
	ctx.Conn.Write(encodedVal)
	return true
}

func Echo(ctx *Context, cmd Command) ReturnValue {
	output := strings.Join(cmd.Args, " ")
	return ReturnValue{RBulkString, output}
}

func Get(ctx *Context, cmd Command) ReturnValue {
	key := cmd.Args[0]
	value, ok := (ctx.State.VariableMap)[key]
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
			return ReturnValue{RBulkString, value.Value}
		}
	}
	return ReturnValue{RNullBulkString, nil}
}

func Incr(ctx *Context, cmd Command) ReturnValue {
	key := cmd.Args[0]
	strVal, ok := (ctx.State.VariableMap)[key]
	var (
		val int
		err error
	)
	if ok {
		val, err = strconv.Atoi(strVal.Value)
		if err != nil {
			return ReturnValue{RSimpleError, ErrorIncr}
		}
	} // otherwise the value is starts at 0
	val += 1
	strVal.Value = strconv.Itoa(val)
	ctx.State.VariableMap[key] = strVal
	return ReturnValue{RInteger, val}
}

func Ping(ctx *Context, cmd Command) ReturnValue {
	return ReturnValue{RSimpleString, "PONG"}
}

func Set(ctx *Context, cmd Command) ReturnValue {
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

	(ctx.State.VariableMap)[key] = Variable{
		Value:              value,
		SetAt:              time.Now().UnixMilli(),
		ExpiryMilliseconds: expiryMilliseconds,
	}
	return ReturnValue{RSimpleString, "OK"}
}

func Type(ctx *Context, cmd Command) ReturnValue {
	key := cmd.Args[0]
	// string, list, set, zset, hash, stream, vectorset

	varMap := ctx.State.VariableMap
	if _, ok := varMap[key]; ok {
		return ReturnValue{
			RSimpleString,
			"string",
		}
	}

	lMap := ctx.State.ListMap
	if _, ok := lMap[key]; ok {
		return ReturnValue{
			RSimpleString,
			"list",
		}
	}

	streamMap := ctx.State.StreamMap
	if _, ok := streamMap[key]; ok {
		return ReturnValue{
			RSimpleString,
			"stream",
		}
	}

	return ReturnValue{
		RSimpleString,
		"none",
	}
}

var CmdFuncMap = map[string]func(ctx *Context, cmd Command) ReturnValue{
	"BLPOP":   Blpop,
	"DISCARD": Discard,
	"ECHO":    Echo,
	"EXEC":    Exec,
	"GET":     Get,
	"INCR":    Incr,
	"LRANGE":  Lrange,
	"PING":    Ping,
	"LLEN":    Llen,
	"LPOP":    Lpop,
	"LPUSH":   Lpush,
	"RPUSH":   Rpush,
	"SET":     Set,
	"TYPE":    Type,
	"XADD":    Xadd,
	"XRANGE":  XRange,
	"XREAD":   XRead,
}
