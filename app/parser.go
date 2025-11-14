package main

import (
	"strconv"
	"strings"
)

type Command struct {
	Command string
	Args    []string
}

func ParseInput(input string) (Command, error) {
	parts := strings.Split(input, "\r\n")
	argCountString := parts[0]
	argCountString = strings.TrimLeft(argCountString, "*")
	argCount, _ := strconv.Atoi(argCountString)
	var command string

	var args []string
	for i := range argCount {
		index := (i * 2) + 1
		_ = parts[index] // verificationSym
		arg := parts[index+1]
		if i == 0 {
			command = arg
			continue
		}
		args = append(args, arg)
	}

	return Command{
		Command: command,
		Args:    args,
	}, nil
}

func ParseXReadArgs(args []string) map[string]string {
	// a,b,c,1,2,3
	// 0,1,2,3,4,5
	size := len(args) - 1
	startStreams := make(map[string]string)
	half := size / 2
	for tIdx := range half {
		keyIdx := tIdx + 1
		valIdx := keyIdx + half
		key := args[keyIdx]
		val := args[valIdx]
		startStreams[key] = val
	}
	return startStreams
}
