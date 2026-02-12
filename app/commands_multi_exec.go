package main

import (
	"strings"
)

func Discard(ctx *Context, cmd Command) ReturnValue {
	if !ctx.ConnState.IsMulti {
		return ReturnValue{
			RSimpleError,
			ErrorDiscardNoMulti,
		}
	}
	ctx.ConnState.IsMulti = false
	return ReturnValue{
		RSimpleString,
		"OK",
	}
}

func Exec(ctx *Context, cmd Command) ReturnValue {
	if !ctx.ConnState.IsMulti {
		return ReturnValue{
			RSimpleError,
			ErrorMultiWithoutExec,
		}
	}
	returnVals := make([]any, 0)
	// execute each of the commands
	for _, cmd := range ctx.ConnState.MultiCmds {
		ret := cmd.Callable(ctx, Command{Args: cmd.Args})
		v := ret.Encoder(ret.EncoderArgs)
		returnVals = append(returnVals, v)
	}

	ctx.ConnState.IsMulti = false
	return ReturnValue{
		RArray,
		returnVals,
	}
}

func Multi(ctx *Context, cmd Command) ReturnValue {
	if !ctx.ConnState.IsMulti {
		ctx.ConnState.IsMulti = true
		return ReturnValue{
			RSimpleString,
			"OK",
		}
	}
	cmdFunc, _ := CmdFuncMap[strings.ToUpper(cmd.Command)]
	ctx.ConnState.MultiCmds = append(ctx.ConnState.MultiCmds, MultiCmd{
		Callable: cmdFunc,
		Args:     cmd.Args,
	})
	return ReturnValue{
		RSimpleString,
		"QUEUED",
	}
}
