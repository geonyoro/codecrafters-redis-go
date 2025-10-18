package main

import "fmt"

func SimpleString(val string) []byte {
	retString := fmt.Sprintf("+%s\r\n", val)
	return []byte(retString)
}

func BulkString(val string) []byte {
	valSize := len(val)
	finalOutput := fmt.Sprintf("$%d\r\n%s\r\n", valSize, val)
	return []byte(finalOutput)
}
