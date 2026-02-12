package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func Xadd(ctx *Context, cmd Command) ReturnValue {
	ctx.State.Mu.Lock()
	defer ctx.State.Mu.Unlock()

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
	ctx.State.Mu.RLock()
	defer ctx.State.Mu.RUnlock()

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
	ctx.State.Mu.RLock()
	defer ctx.State.Mu.RUnlock()

	args, err := ParseXReadArgs(cmd.Args)
	if err != nil {
		return ReturnValue{
			RSimpleError,
			"ERR Invalid argumen to Xread",
		}
	}
	retArray := make([]any, 0)
	rets, isNilArray := xReadInner(ctx, args)
	if isNilArray {
		return ReturnValue{RNullArray, 1}
	}
	for _, ret := range rets {
		retArray = append(retArray, ret.ToRArray())
	}

	return ReturnValue{
		RArray,
		retArray,
	}
}

type XReadArgs struct {
	Streams    map[string]string
	Count      int
	Block      int
	IsBlocking bool
}

func ParseXReadArgs(args []string) (XReadArgs, error) {
	// count,2,block,100,streams,a,b,c,1,2,3
	// block,100,streams,0,1,2,3,4,5
	// read count and block section

	size := len(args)
	var count int
	var block int
	isBlocking := false
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
				isBlocking = true
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
		Streams:    streams,
		Count:      count,
		Block:      block,
		IsBlocking: isBlocking,
	}, nil
}

func xReadInner(ctx *Context, args XReadArgs) (ret []XReadReturn, isNilArray bool) {
	for streamId, fromId := range args.Streams {
		streamEntries := make([]XRangeReturn, 0)
		stream := ctx.State.StreamMap[streamId]
		fromId, err := incFromId(fromId, stream)
		if err != nil {
			panic(err)
		}
		xrangeRets := xRangeInner(ctx, streamId, fromId, "+")
		if args.IsBlocking && len(xrangeRets) <= 0 {
			// poll the stream every Block milliseconds
			var timerExpiry *time.Timer
			if args.Block > 0 {
				timerExpiry = time.NewTimer(time.Duration(args.Block) * time.Millisecond)
			}
			// 10ms ticker
			ticker := time.NewTicker(10 * time.Millisecond)

		CheckLoop:
			for {
				select {
				case <-ticker.C:
					xrangeRets = xRangeInner(ctx, streamId, fromId, "+")
					if len(xrangeRets) > 0 {
						ticker.Stop()
						break CheckLoop
					}
				default:
				}

				if args.Block > 0 {
					select {
					case <-timerExpiry.C:
						// the timer has expired
						ticker.Stop()
						break CheckLoop
					default:
					}
				}
			}
			if len(xrangeRets) == 0 {
				// we have to respond with a nil array
				return []XReadReturn{}, true
			}
		}
		for _, xrangeRet := range xrangeRets {
			streamEntries = append(streamEntries, xrangeRet)
		}
		streamRet := XReadReturn{
			streamId,
			streamEntries,
		}
		ret = append(ret, streamRet)
	}
	return ret, false
}

func xRangeInner(ctx *Context, streamId, fromId, toId string) (ret []XRangeReturn) {
	ret = make([]XRangeReturn, 0)
	stream, ok := ctx.State.StreamMap[streamId]
	if !ok {
		return ret
	}

	// convert the from section
	var fromMillis, fromSequence int
	if fromId == "-" {
		fromMillis = 0
	} else {
		fromMillis, fromSequence = toMillisSeq(fromId)
	}

	// convert the to section
	var toMillis, toSequence int
	if toId == "+" {
		toMillis, toSequence = stream.GetLastMillisSequence()
	} else {
		toMillis, toSequence = toMillisSeq(toId)
	}

	for _, keyStr := range stream.Keys {
		key, _ := strconv.Atoi(keyStr)
		// do some bounds checking
		if key > toMillis {
			// guard, we have exceeded the millis
			break
		}
		if key < fromMillis {
			// guard, we have not reached the millis yet
			continue
		}

		// we are within the valid millis bounds

		// check we have reached the fromMillis Sequence
		if key == fromMillis {
			millisVal := stream.Map[keyStr]
			for _, sequenceId := range millisVal.Keys {
				sequenceInt, _ := strconv.Atoi(sequenceId)
				if sequenceInt < fromSequence {
					// guard: we have not yet reached the sequence yet
					continue
				}
				if key == toMillis && sequenceInt > toSequence {
					// very specific condition where the 2 are equal
					// guard: we have exceeded the sequence
					break
				}
				ret = append(ret, XRangeReturn{
					ID: fmt.Sprintf("%d-%d", key, sequenceInt),
					KV: millisVal.Map[sequenceId],
				})
			}
			continue
		}

		// check we have not exceed the toMillis Sequence
		if key == toMillis {
			millisVal := stream.Map[keyStr]
			for _, sequenceId := range millisVal.Keys {
				sequenceInt, _ := strconv.Atoi(sequenceId)
				if sequenceInt > toSequence {
					// guard: we have exceeded the guard bounds
					continue
				}
				ret = append(ret, XRangeReturn{
					ID: fmt.Sprintf("%d-%d", key, sequenceInt),
					KV: millisVal.Map[sequenceId],
				})
			}
			continue
		}

		// just add everything in between here
		millisVal := stream.Map[keyStr]
		for _, sequenceId := range millisVal.Keys {
			ret = append(ret, XRangeReturn{
				ID: fmt.Sprintf("%s-%s", keyStr, sequenceId),
				KV: millisVal.Map[sequenceId],
			})
		}
	}

	return ret
}

func incFromId(fromId string, stream *Stream) (string, error) {
	var millis, sequence int
	if fromId == "$" {
		if stream == nil {
			return "0-0", nil
		}
		millis, sequence = stream.GetLastMillisSequence()
	} else {
		millis, sequence = toMillisSeq(fromId)
	}
	return fmt.Sprintf("%d-%d", millis, sequence+1), nil
}

func toMillisSeq(id string) (int, int) {
	Parts := strings.Split(id, "-")
	var millis, sequence int
	millis, _ = strconv.Atoi(Parts[0])
	if len(Parts) > 1 {
		sequence, _ = strconv.Atoi(Parts[1])
	}
	return millis, sequence
}
