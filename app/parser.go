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
