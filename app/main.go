package main

import (
	"fmt"
	"log/slog"
	"net"
	"os"
)

func main() {
	globalState := NewState()
	args := ParseCliArgs(os.Args[1:])
	globalState.updateWithCliArgs(args)
	masterConn, err := SetupReplication(globalState)
	if err != nil {
		panic(err)
	}

	if globalState.IsReplica() {
		go handleConn(*masterConn, globalState, true)
	}

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
		go handleConn(conn, globalState, false)
	}
}

func handleConn(conn net.Conn, globalState *State, inMasterConn bool) {
	defer conn.Close()
	connState := NewConnState()

	buffer := make([]byte, 1024)
	for {
		readSize, err := conn.Read(buffer)
		if err != nil || readSize == 0 {
			break
		}
		commands, err := ParseInput(buffer[:readSize])
		if err != nil {
			fmt.Println(err)
			continue
		}
		for _, command := range commands {
			fmt.Printf("%v inMasterConn:%t %+v %s\n", command, inMasterConn, globalState.VariableMap, buffer)

			if IsWriteCommand(command.Command) && len(globalState.Settings.Replicas) > 0 {
				for _, replica := range globalState.Settings.Replicas {
					(*replica).Write(buffer[:readSize])
				}
			}

			ctx := Context{
				Conn:      conn,
				ConnState: connState,
				State:     globalState,
			}
			ExecuteCommand(&ctx, command)
		}
	}
}
