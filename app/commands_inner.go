package main

import (
	"fmt"
	"strconv"
	"strings"
)

func xRangeInner(ctx *Context, streamId, fromId, toId string) (ret []XRangeReturn) {
	ret = make([]XRangeReturn, 0)
	stream, ok := ctx.State.StreamMap[streamId]
	if !ok {
		return ret
	}

	// convert the from section
	var fromMillis, fromSequence int
	fromParts := strings.Split(fromId, "-")
	fromMillis, _ = strconv.Atoi(fromParts[0])
	if len(fromParts) > 1 {
		fromSequence, _ = strconv.Atoi(fromParts[1])
	}

	// convert the to section
	var toMillis, toSequence int
	toParts := strings.Split(toId, "-")
	toMillis, _ = strconv.Atoi(toParts[0])
	if len(toParts) > 1 {
		toSequence, _ = strconv.Atoi(toParts[1])
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
