package main

import (
	"fmt"
)

func RSimpleString(arg any) []byte {
	// The string mustn't contain a CR (\r) or LF (\n) character and is terminated by CRLF (i.e., \r\n).
	val := arg.(string)
	retString := fmt.Sprintf("+%s\r\n", val)
	return []byte(retString)
}

func RSimpleError(arg any) []byte {
	val := arg.(string)
	retString := fmt.Sprintf("-%s\r\n", val)
	return []byte(retString)
}

func RBulkString(arg any) []byte {
	// A bulk string represents a single binary string.
	val := arg.(string)
	valSize := len(val)
	output := fmt.Sprintf("$%d\r\n%s\r\n", valSize, val)
	return []byte(output)
}

func RNullBulkString(arg any) []byte {
	return []byte("$-1\r\n")
}

func RNullArray(arg any) []byte {
	return []byte("*-1\r\n")
}

func RInteger(arg any) []byte {
	size := arg.(int)
	sign := ""
	if size < 0 {
		sign = "-"
	}
	output := fmt.Sprintf(":%s%d\r\n", sign, size)
	return []byte(output)
}

func RArray(arg any) []byte {
	elems := arg.([]any)
	size := len(elems)
	outputStr := fmt.Sprintf("*%d\r\n", size)
	output := []byte(outputStr)
	for _, elem := range elems {
		switch e := elem.(type) {
		case string:
			output = append(output, RBulkString(e)...)
		case int:
			output = append(output, RInteger(e)...)
		case []any:
			ret := RArray(e)
			fmt.Printf("%q\n", string(ret))
			output = append(output, ret...)
		default:
			panic(fmt.Sprintf("Unknown type: %T", elem))
		}
	}
	fmt.Printf("%q\n", string(output))
	return []byte(output)
}
