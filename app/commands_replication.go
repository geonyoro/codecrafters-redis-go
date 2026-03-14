package main

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net"
	"strings"
)

var InvalidReplica = errors.New("invalid replica")

func SetupReplication(globalState *State) (conn *net.Conn, err error) {
	if globalState.Settings.ReplicaOf == "" {
		return
	}

	parts := strings.Split(globalState.Settings.ReplicaOf, " ")
	if len(parts) != 2 {
		err = InvalidReplica
		return
	}
	host := parts[0]
	port := parts[1]

	var c net.Conn
	c, err = Connect(host, port)
	conn = &c
	if err != nil {
		return
	}

	err = ReplPing(c)
	if err != nil {
		return
	}

	err = ReplConf(c, globalState)
	if err != nil {
		return
	}

	err = ReplPsync(c, globalState)
	if err != nil {
		return
	}
	fmt.Print("SetupReplication finished.\n")
	return
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
	fmt.Print("Ping command finished.\n")
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
	fmt.Print("ReplConfPort finished.\n")
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
	fmt.Print("ReplConfCapa finished.\n")
	return nil
}

func ReplPsync(conn net.Conn, globalState *State) error {
	val := RArray([]any{"PSYNC", "?", "-1"})
	_, err := conn.Write(val)
	if err != nil {
		return err
	}

	buffer := make([]byte, 1024)
	_, err = conn.Read(buffer)
	if err != nil {
		return err
	}
	fmt.Print("ReplPsync finished.\n")
	return nil
}

func ReplConfAsMaster(ctx *Context, cmd Command) ReturnValue {
	return ReturnValue{RSimpleString, "OK"}
}

func PsyncAsMaster(ctx *Context, cmd Command) ReturnValue {
	// the master cannot perform an incremental update to the replica, and will start a full resynchronization
	ctx.Conn.Write(RSimpleString(fmt.Sprintf("FULLRESYNC %s 0", ctx.State.Settings.MasterReplId)))

	rdb64 := "UkVESVMwMDEx+glyZWRpcy12ZXIFNy4yLjD6CnJlZGlzLWJpdHPAQPoFY3RpbWXCbQi8ZfoIdXNlZC1tZW3CsMQQAPoIYW9mLWJhc2XAAP/wbjv+wP9aog=="
	rdb, err := base64.StdEncoding.DecodeString(rdb64)
	if err != nil {
		panic(err)
	}
	ctx.Conn.Write(RRawBytes(rdb))

	ctx.State.Settings.Replicas = append(ctx.State.Settings.Replicas, &ctx.Conn)
	return ReturnValue{REmpty, ""}
}
