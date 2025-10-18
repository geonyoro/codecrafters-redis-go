package main

import (
	"fmt"
	"log/slog"
	"net"
	"os"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	globalState := NewState()

	l, err := net.Listen("tcp4", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	slog.Debug("Bound: Accepting connections.")
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
			Conn:  conn,
			State: globalState,
		}
		ExecuteCommand(&ctx, command)
	}
}
