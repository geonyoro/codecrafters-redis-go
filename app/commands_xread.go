package main

import (
	"fmt"
	"strconv"
	"strings"
)

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
