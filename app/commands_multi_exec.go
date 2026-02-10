package main

import "strings"

func Exec(ctx *Context, cmd Command) ReturnValue {
	if !ctx.State.IsMulti {
		return ReturnValue{
			RSimpleError,
			ErrorMultiWithoutExec,
		}
	}
	returnVals := make([]any, 0)
	// execute each of the commands
	for _, cmd := range ctx.State.MultiCmds {
		ret := cmd.Callable(ctx, Command{Args: cmd.Args})
		returnVals = append(returnVals, ret)
	}

	ctx.State.IsMulti = false
	return ReturnValue{
		RArray,
		returnVals,
	}
}

func Multi(ctx *Context, cmd Command) ReturnValue {
	if !ctx.State.IsMulti {
		ctx.State.IsMulti = true
		return ReturnValue{
			RSimpleString,
			"OK",
		}
	}
	cmdFunc, _ := CmdFuncMap[strings.ToUpper(cmd.Command)]
	ctx.State.MultiCmds = append(ctx.State.MultiCmds, MultiCmd{
		Callable: cmdFunc,
		Args:     cmd.Args,
	})
	return ReturnValue{
		RSimpleString,
		"QUEUED",
	}
}
