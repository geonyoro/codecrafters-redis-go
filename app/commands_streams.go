package main

import (
	"fmt"
	"strconv"
	"strings"
)

func Xadd(ctx *Context, cmd Command) ReturnValue {
	streamKey := cmd.Args[0]
	stream := ctx.State.GetOrCreateStreamForKey(streamKey)

	id := cmd.Args[1]
	if id == "0-0" {
		return ReturnValue{
			RSimpleError,
			"ERR The ID specified in XADD must be greater than 0-0",
		}
	}

	var millis, sequence string
	if id == "*" {
		// full generation mode
		millis = stream.GenerateMillis()
		sequence = stream.GenerateSequence(millis)
	} else {
		stringParts := strings.Split(id, "-")
		millis = stringParts[0]
		sequence = stringParts[1]
		isValid, err := stream.IsNewStreamIdValid(millis, sequence)
		if err != nil {
			return ReturnValue{
				RSimpleError,
				"ERR Unknown Error",
			}
		}
		if !isValid {
			return ReturnValue{
				RSimpleError,
				"ERR The ID specified in XADD is equal or smaller than the target stream top item",
			}
		}

		if sequence == "*" {
			sequence = stream.GenerateSequence(millis)
		}
	}

	// first 2 entries are taken by stream and id
	// iterate in batches of 2
	kvMap := make(map[string]string)
	for i := range (len(cmd.Args) - 2) / 2 {
		idx := (i + 1) * 2
		key, value := cmd.Args[idx], cmd.Args[idx+1]
		kvMap[key] = value
	}
	err := stream.AddIdWithKV(millis, sequence, kvMap)
	if err != nil {
		return ReturnValue{
			RSimpleError,
			"ERR Unknown Error",
		}
	}
	id = fmt.Sprintf("%s-%s", millis, sequence)
	return ReturnValue{
		RBulkString,
		id,
	}
}

func XRange(ctx *Context, cmd Command) ReturnValue {
	streamKey := cmd.Args[0]
	fromId := cmd.Args[1]
	toId := cmd.Args[2]

	retArray := make([]any, 0)
	rets := xRangeInner(ctx, streamKey, fromId, toId)
	for _, ret := range rets {
		retArray = append(retArray, ret.ToRArray())
	}

	return ReturnValue{
		RArray,
		retArray,
	}
}

func XRead(ctx *Context, cmd Command) ReturnValue {
	args, err := ParseXReadArgs(cmd.Args)
	if err != nil {
		return ReturnValue{
			RSimpleError,
			"ERR Invalid argumen to Xread",
		}
	}
	retArray := make([]any, 0)
	rets := xReadInner(ctx, args.Streams)
	for _, ret := range rets {
		retArray = append(retArray, ret.ToRArray())
	}

	return ReturnValue{
		RArray,
		retArray,
	}
}

type XReadArgs struct {
	Streams map[string]string
	Count   int
	Block   int
}

func ParseXReadArgs(args []string) (XReadArgs, error) {
	// count,2,block,100,streams,a,b,c,1,2,3
	// block,100,streams,0,1,2,3,4,5
	// read count and block section

	size := len(args)
	var count int
	var block int
	isStreamLiteralSeen := false
	streamStartsAt := 0
	for i := 0; ; i++ { // go until you hit the streams section
		key := args[i]
		if strings.ToLower(key) == "block" || strings.ToLower(key) == "count" {
			i += 1
			valStr := args[i]
			val, err := strconv.Atoi(valStr)
			if err != nil {
				return XReadArgs{}, err
			}
			if strings.ToLower(key) == "count" {
				count = val
			} else {
				block = val
			}
			continue
		}
		// we have hit streams
		if !isStreamLiteralSeen {
			if strings.ToLower(key) != "streams" {
				return XReadArgs{}, fmt.Errorf("expected streams section to start with 'STREAMS'")
			}
			isStreamLiteralSeen = true
			streamStartsAt = i + 1
			continue
		}
		break
	}
	streams := make(map[string]string)
	streamSectionSize := size - streamStartsAt
	half := streamSectionSize/2 + streamStartsAt
	for i := 0; i < streamSectionSize/2; i += 1 {
		key := args[streamStartsAt+i]
		val := args[half+i]
		streams[key] = val
	}

	return XReadArgs{
		Streams: streams,
		Count:   count,
		Block:   block,
	}, nil
}
