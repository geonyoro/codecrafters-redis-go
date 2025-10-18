package main

import "fmt"

func RSimpleString(val string) []byte {
	retString := fmt.Sprintf("+%s\r\n", val)
	return []byte(retString)
}

func RBulkString(val string) []byte {
	valSize := len(val)
	finalOutput := fmt.Sprintf("$%d\r\n%s\r\n", valSize, val)
	return []byte(finalOutput)
}

func RNullBulkString() []byte {
	return []byte("$-1\r\n")
}

func RInteger(size int) []byte {
	sign := ""
	if size < 0 {
		sign = "-"
	}
	output := fmt.Sprintf(":%s%d\r\n", sign, size)
	fmt.Println(output)
	return []byte(output)
}
