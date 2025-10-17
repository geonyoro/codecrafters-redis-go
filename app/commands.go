package main

import (
	"fmt"
	"net"
	"strings"
)

func ExecuteCommand(conn net.Conn, cmd Command) {
	if cmd.Command == "ECHO" {
		Echo(conn, cmd.Args)
	}
}

func Echo(conn net.Conn, args []string) {
	output := strings.Join(args, " ")
	outputSize := len(output)
	finalOutput := fmt.Sprintf("$%d\r\n%s\r\n", outputSize, output)
	conn.Write([]byte(finalOutput))
}
