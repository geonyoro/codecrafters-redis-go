package main

func Exec(ctx *Context, cmd Command) ReturnValue {
	if !ctx.State.IsMulti {
		return ReturnValue{
			RSimpleError,
			ErrorMultiWithoutExec,
		}
	}

	ctx.State.IsMulti = false
	returnVals := make([]any, 0)
	if len(ctx.State.MultiCmds) == 0 {
		return ReturnValue{
			RArray,
			returnVals,
		}
	}
	return ReturnValue{
		RSimpleString,
		"OK",
	}
}

func Multi(ctx *Context, cmd Command) ReturnValue {
	ctx.State.IsMulti = true
	return ReturnValue{
		RSimpleString,
		"OK",
	}
}
