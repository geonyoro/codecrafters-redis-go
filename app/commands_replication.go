package main

import (
	"errors"
	"fmt"
	"net"
	"strings"
)

var InvalidReplica = errors.New("invalid replica")

func SetupReplication(globalState *State) error {
	if globalState.Settings.ReplicaOf == "" {
		return nil
	}

	parts := strings.Split(globalState.Settings.ReplicaOf, " ")
	if len(parts) != 2 {
		return InvalidReplica
	}
	host := parts[0]
	port := parts[1]

	conn, err := Connect(host, port)
	if err != nil {
		return err
	}

	err = ReplPing(conn)
	if err != nil {
		return err
	}

	err = ReplConf(conn, globalState)
	if err != nil {
		return err
	}

	return nil
}

func Connect(host, port string) (net.Conn, error) {
	return net.Dial("tcp", fmt.Sprintf("%s:%s", host, port))
}

func ReplPing(conn net.Conn) error {
	val := RArray([]any{"Ping"})
	_, err := conn.Write(val)
	if err != nil {
		return err
	}

	// await response
	buffer := make([]byte, 1024)
	_, err = conn.Read(buffer)
	if err != nil {
		return err
	}
	return nil
}

func ReplConf(conn net.Conn, globalState *State) error {
	var err error
	err = ReplConfPort(conn, globalState)
	if err != nil {
		return err
	}

	err = ReplConfCapa(conn, globalState)
	if err != nil {
		return err
	}

	return nil
}

func ReplConfPort(conn net.Conn, globalState *State) error {
	val := RArray([]any{"REPLCONF", "listening-port", fmt.Sprintf("%d", globalState.Settings.Port)})
	_, err := conn.Write(val)
	if err != nil {
		return err
	}

	buffer := make([]byte, 1024)
	_, err = conn.Read(buffer)
	if err != nil {
		return err
	}
	return nil
}

func ReplConfCapa(conn net.Conn, globalState *State) error {
	val := RArray([]any{"REPLCONF", "capa", "psync2"})
	_, err := conn.Write(val)
	if err != nil {
		return err
	}

	buffer := make([]byte, 1024)
	_, err = conn.Read(buffer)
	if err != nil {
		return err
	}
	return nil
}
