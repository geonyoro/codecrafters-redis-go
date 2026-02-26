package main

import (
	"log"
	"os"
	"strconv"
)

type CliArgs struct {
	Port int
	Host string
}

func parseCliArgs() *CliArgs {
	c := &CliArgs{
		Port: 6379,
		Host: "127.0.0.1",
	}

	i := 0
	size := len(os.Args)
	for {
		i += 1
		if i >= size {
			break
		}

		if os.Args[i] == "-p" || os.Args[i] == "--port" {
			i += 1
			port, err := strconv.Atoi(os.Args[i])
			if err != nil {
				log.Fatal("invalid argument: --port/-p must be number")
			}
			c.WithPort(port)
		}

		if os.Args[i] == "-h" || os.Args[i] == "--host" {
			i += 1
			c.WithHost(os.Args[i])
		}
	}
	return c
}

func (c *CliArgs) WithPort(port int) *CliArgs {
	c.Port = port
	return c
}

func (c *CliArgs) WithHost(host string) *CliArgs {
	c.Host = host
	return c
}
