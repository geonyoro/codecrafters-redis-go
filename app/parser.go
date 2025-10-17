package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Command struct {
	Command string
	Args    []string
}

func ParseInput(input string) (Command, error) {
	parts := strings.Split(input, "\r\n")
	fmt.Printf("%+v\n", parts)
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
