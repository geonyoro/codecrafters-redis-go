package main

import "fmt"

func Info(ctx *Context, cmd Command) ReturnValue {
	role := "master"
	if len(ctx.State.Settings.ReplicaOf) > 0 {
		role = "slave"
	}
	return ReturnValue{
		RBulkString, fmt.Sprintf(`
role:%s
connected_slaves:0
master_replid:%s
master_repl_offset:%d
second_repl_offset:-1
repl_backlog_active:0
repl_backlog_size:1048576
repl_backlog_first_byte_offset:0
repl_backlog_histlen:`, role, ctx.State.Settings.MasterReplId, ctx.State.Settings.MasterReplOffset),
	}
}
