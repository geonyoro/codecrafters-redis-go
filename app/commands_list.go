package main

import (
	"strconv"
	"time"
)

func Blpop(ctx *Context, cmd Command) ReturnValue {
	ctx.State.Mu.Lock()
	defer ctx.State.Mu.Unlock()

	listName := cmd.Args[0]
	listMap := ctx.State.ListMap
	list, ok := listMap[listName]
	if !ok {
		// create the list
		list = &ListVariable{}
		listMap[listName] = list
		ctx.State.ListMap = listMap
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
				EncoderArgs: []any{listName, elem},
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

func Llen(ctx *Context, cmd Command) ReturnValue {
	ctx.State.Mu.RLock()
	defer ctx.State.Mu.RUnlock()

	listName := cmd.Args[0]
	listMap := ctx.State.ListMap
	list, ok := listMap[listName]
	if !ok {
		return ReturnValue{RInteger, 0}
	}
	return ReturnValue{RInteger, len(list.Values)}
}

func Lpop(ctx *Context, cmd Command) ReturnValue {
	ctx.State.Mu.Lock()
	defer ctx.State.Mu.Unlock()

	listName := cmd.Args[0]
	listMap := ctx.State.ListMap
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
		return ReturnValue{RBulkString, value}
	}
	size, _ := strconv.Atoi(cmd.Args[1])
	values := list.Values[:size]
	list.Values = list.Values[size:]
	return ReturnValue{RArray, StringArraytoAny(values)}
}

func Lpush(ctx *Context, cmd Command) ReturnValue {
	ctx.State.Mu.Lock()
	defer ctx.State.Mu.Unlock()

	listName := cmd.Args[0]
	listMap := ctx.State.ListMap
	list, ok := listMap[listName]
	if !ok {
		list = &ListVariable{}
		listMap[listName] = list
		ctx.State.ListMap = listMap
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
	ctx.State.Mu.Lock()
	defer ctx.State.Mu.Unlock()

	listName := cmd.Args[0]
	listMap := ctx.State.ListMap
	list, ok := listMap[listName]
	if !ok {
		list = &ListVariable{}
		listMap[listName] = list
		ctx.State.ListMap = listMap
	}
	for i := 1; i < len(cmd.Args); i++ {
		newValue := cmd.Args[i]
		list.Values = append(list.Values, newValue)
	}
	return ReturnValue{RInteger, len(list.Values)}
}

func Lrange(ctx *Context, cmd Command) ReturnValue {
	ctx.State.Mu.RLock()
	defer ctx.State.Mu.RUnlock()

	listName := cmd.Args[0]
	listMap := ctx.State.ListMap
	listVar, ok := listMap[listName]
	if !ok {
		return ReturnValue{RArray, []any{}}
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
		output := []any{}
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
	return ReturnValue{RArray, StringArraytoAny(values)}
}
