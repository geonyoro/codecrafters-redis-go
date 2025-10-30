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
	cmdFunc, ok := CmdFuncMap[strings.ToUpper(cmd.Command)]
	if ok {
		returnVal := cmdFunc(ctx, cmd)
		encodedVal := returnVal.Encoder(returnVal.EncoderArgs)
		ctx.Conn.Write(encodedVal)
		return true
	}
	fmt.Println("Failed to find cmd for", cmd.Command)
	return false
}

func Blpop(ctx *Context, cmd Command) ReturnValue {
	listName := cmd.Args[0]
	listMap := *ctx.State.ListMap
	list, ok := listMap[listName]
	if !ok {
		// create the list
		list = &ListVariable{}
		listMap[listName] = list
	}

	// access timeout
	var timeout float64 = 0
	if len(cmd.Args) > 1 {
		if ttimeout, err := strconv.ParseFloat(cmd.Args[1], 64); err == nil {
			timeout = ttimeout * 1000
		}
	}
	startTime := time.Now()
	var endTime time.Time
	if timeout > 0 {
		endTime = startTime.Add(time.Millisecond * time.Duration(timeout))
	}
	for {
		if len(list.Values) > 0 {
			elem := list.Values[0]
			list.Values = list.Values[1:len(list.Values)]
			return ReturnValue{
				Encoder:     RArray,
				EncoderArgs: []string{listName, elem},
			}
		}
		if timeout > 0 {
			if time.Now().After(endTime) {
				// it has expired
				return ReturnValue{
					Encoder:     RNullArray,
					EncoderArgs: 1,
				}
			}
		}
		time.Sleep(time.Millisecond * time.Duration(10))
	}
}

func Echo(ctx *Context, cmd Command) ReturnValue {
	output := strings.Join(cmd.Args, " ")
	return ReturnValue{RBulkString, output}
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

	(*ctx.State.VariableMap)[key] = Variable{
		Value:              value,
		SetAt:              time.Now().UnixMilli(),
		ExpiryMilliseconds: expiryMilliseconds,
	}
	return ReturnValue{RSimpleString, "OK"}
}

func Get(ctx *Context, cmd Command) ReturnValue {
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
			return ReturnValue{RBulkString, value.Value}
		}
	}
	return ReturnValue{RNullBulkString, nil}
}

func Llen(ctx *Context, cmd Command) ReturnValue {
	listName := cmd.Args[0]
	listMap := *ctx.State.ListMap
	list, ok := listMap[listName]
	if !ok {
		return ReturnValue{RInteger, 0}
	}
	return ReturnValue{RInteger, len(list.Values)}
}

func Lpop(ctx *Context, cmd Command) ReturnValue {
	listName := cmd.Args[0]
	listMap := *ctx.State.ListMap
	list, ok := listMap[listName]
	if !ok || len(list.Values) == 0 {
		return ReturnValue{
			Encoder:     RNullBulkString,
			EncoderArgs: 1,
		}
	}
	if len(cmd.Args) == 1 {
		// remove a single value
		value := list.Values[0]
		list.Values = list.Values[1:len(list.Values)]
		return ReturnValue{RSimpleString, value}
	}
	size, _ := strconv.Atoi(cmd.Args[1])
	values := list.Values[:size]
	list.Values = list.Values[size:]
	return ReturnValue{RArray, values}
}

func Lpush(ctx *Context, cmd Command) ReturnValue {
	listName := cmd.Args[0]
	listMap := *ctx.State.ListMap
	list, ok := listMap[listName]
	if !ok {
		list = &ListVariable{}
		listMap[listName] = list
	}
	argSize := len(cmd.Args) - 1
	newList := make([]string, argSize)
	for i, arg := range cmd.Args {
		if i == 0 {
			continue
		}
		newList[argSize-i] = arg
	}
	newList = append(newList, list.Values...)
	list.Values = newList
	return ReturnValue{RInteger, len(list.Values)}
}

func Rpush(ctx *Context, cmd Command) ReturnValue {
	listName := cmd.Args[0]
	listMap := *ctx.State.ListMap
	list, ok := listMap[listName]
	if !ok {
		list = &ListVariable{}
		listMap[listName] = list
	}
	for i := 1; i < len(cmd.Args); i++ {
		newValue := cmd.Args[i]
		list.Values = append(list.Values, newValue)
	}
	return ReturnValue{RInteger, len(list.Values)}
}

func convertLrangeIndex(listSize int, index int) uint {
	maxNegativeIndex := 0 - listSize
	// manipulate the indexes
	if index < 0 {
		if index < maxNegativeIndex {
			index = 0
		} else {
			index = listSize + index
		}
	} else {
		if index >= listSize {
			index = listSize - 1
		}
	}
	return uint(index)
}

func Lrange(ctx *Context, cmd Command) ReturnValue {
	listName := cmd.Args[0]
	listMap := *ctx.State.ListMap
	listVar, ok := listMap[listName]
	if !ok {
		return ReturnValue{RArray, []string{}}
	}
	listSize := len(listVar.Values)
	// access the indexes
	startIndex, _ := strconv.Atoi(cmd.Args[1])
	endIndex, _ := strconv.Atoi(cmd.Args[2])
	maxNegativeIndex := 0 - listSize
	// manipulate the indexes
	if startIndex < 0 {
		if startIndex < maxNegativeIndex {
			startIndex = 0
		} else {
			startIndex = len(listVar.Values) + startIndex
		}
	}
	if endIndex < 0 {
		if endIndex < maxNegativeIndex {
			endIndex = 0
		} else {
			endIndex = len(listVar.Values) + endIndex
		}
	}
	if !ok || startIndex >= endIndex {
		output := []string{}
		return ReturnValue{RArray, output}
	}
	if endIndex >= listSize {
		endIndex = listSize - 1
	}
	accessSize := (endIndex - startIndex) + 1
	values := make([]string, accessSize)

	for i := range accessSize {
		values[i] = listVar.Values[i+startIndex]
	}
	return ReturnValue{RArray, values}
}

func Type(ctx *Context, cmd Command) ReturnValue {
	key := cmd.Args[0]
	// string, list, set, zset, hash, stream, vectorset
	state := *ctx.State

	varMap := *state.VariableMap
	if _, ok := varMap[key]; ok {
		return ReturnValue{
			RSimpleString,
			"string",
		}
	}

	lMap := *state.ListMap
	if _, ok := lMap[key]; ok {
		return ReturnValue{
			RSimpleString,
			"list",
		}
	}

	streamMap := *state.StreamMap
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

func Xadd(ctx *Context, cmd Command) ReturnValue {
	streamKey := cmd.Args[0]
	streamMap := *ctx.State.StreamMap
	stream, ok := streamMap[streamKey]
	if !ok {
		// make the stream
		stream = &Stream{
			Entries: make(map[string]Entry),
		}
		streamMap[streamKey] = stream
	}

	id := cmd.Args[1]
	if id == "0-0" {
		return ReturnValue{
			RSimpleError,
			"ERR The ID specified in XADD must be greater than 0-0",
		}
	}

	stringParts := strings.Split(id, "-")
	millis, _ := strconv.Atoi(stringParts[0])
	sequence, _ := strconv.Atoi(stringParts[1])
	if !IsValidNewStreamId(stream.LastEntry, millis, sequence) {
		return ReturnValue{
			RSimpleError,
			"ERR The ID specified in XADD is equal or smaller than the target stream top item",
		}
	}
	stream.LastEntry = []int{millis, sequence}

	stream.Entries[id] = Entry{}
	// first 2 entries are taken by stream and id
	// iterate in batches of 2
	for i := range (len(cmd.Args) - 2) / 2 {
		idx := (i + 1) * 2
		key, value := cmd.Args[idx], cmd.Args[idx+1]
		stream.Entries[id][key] = value
	}
	return ReturnValue{
		RBulkString,
		id,
	}
}

var CmdFuncMap = map[string]func(ctx *Context, cmd Command) ReturnValue{
	"BLPOP":  Blpop,
	"ECHO":   Echo,
	"GET":    Get,
	"LRANGE": Lrange,
	"PING":   Ping,
	"LLEN":   Llen,
	"LPOP":   Lpop,
	"LPUSH":  Lpush,
	"RPUSH":  Rpush,
	"SET":    Set,
	"TYPE":   Type,
	"XADD":   Xadd,
}
