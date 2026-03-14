package main

import (
	"log"
	"strconv"
	"strings"
)

type Command struct {
	Command string
	Args    []string
}

func ParseInput(inputb []byte) ([]Command, error) {
	parts := strings.Split(string(inputb), "\r\n")
	idx := 0
	partCount := len(parts)
	commands := make([]Command, 0)
	for {
		if idx >= partCount-1 {
			break
		}
		// we expect many arrays
		argCountString := parts[idx]
		argCountString = strings.TrimLeft(argCountString, "*")
		argCount, err := strconv.Atoi(argCountString)
		if err != nil {
			return commands, err
		}

		var command string

		var args []string
		for argIdx := range argCount {
			index := (argIdx * 2) + 1 + idx
			_ = parts[index] // verificationSym
			arg := parts[index+1]
			if argIdx == 0 {
				command = arg
				continue
			}
			args = append(args, arg)
		}
		idx = idx + argCount*2 + 1
		commands = append(commands, Command{
			Command: command,
			Args:    args,
		})
	}
	return commands, nil
}

func ParseCliArgs(args []string) *CliArgs {
	c := &CliArgs{
		Port: 6379,
		Host: "127.0.0.1",
	}

	i := -1
	size := len(args)
	for {
		i += 1
		if i >= size {
			break
		}

		if args[i] == "-p" || args[i] == "--port" {
			i += 1
			port, err := strconv.Atoi(args[i])
			if err != nil {
				log.Fatal("invalid argument: --port/-p must be number")
			}
			c.WithPort(port)
		}

		if args[i] == "-h" || args[i] == "--host" {
			i += 1
			c.WithHost(args[i])
		}

		if args[i] == "--replicaof" {
			i += 1
			c.WithReplicaOf(args[i])
		}
	}
	return c
}
