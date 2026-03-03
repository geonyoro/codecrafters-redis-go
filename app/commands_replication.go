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
		return nil
	}
	ReplPing(conn)
	return nil
}

func Connect(host, port string) (net.Conn, error) {
	return net.Dial("tcp", fmt.Sprintf("%s:%s", host, port))
}

func ReplPing(conn net.Conn) error {
	val := RArray([]any{"Ping"})
	fmt.Printf("%s", val)
	_, err := conn.Write(val)
	if err != nil {
		return err
	}
	return nil
}
