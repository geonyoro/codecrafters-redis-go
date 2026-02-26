package main

import (
	"fmt"
	"log/slog"
	"net"
	"os"
)

func main() {
	globalState := NewState()
	args := parseCliArgs()

	l, err := net.Listen("tcp4", fmt.Sprintf("%s:%d", args.Host, args.Port))
	if err != nil {
		fmt.Printf("Failed to bind to host:port %s:%d\n", args.Host, args.Port)
		os.Exit(1)
	}

	slog.Debug(fmt.Sprintf("Bound: Accepting connections on %s:%d.", args.Host, args.Port))
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			continue
		}
		go handleConn(conn, globalState)
	}
}

func handleConn(conn net.Conn, globalState *State) {
	defer conn.Close()
	connState := NewConnState()

	for {
		buffer := make([]byte, 1024)
		readSize, err := conn.Read(buffer)
		if err != nil || readSize == 0 {
			break
		}
		command, err := ParseInput(string(buffer))
		if err != nil {
			fmt.Println(err)
			continue
		}
		ctx := Context{
			Conn:      conn,
			ConnState: connState,
			State:     globalState,
		}
		ExecuteCommand(&ctx, command)
	}
}
