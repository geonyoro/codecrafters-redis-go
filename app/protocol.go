package main

import "fmt"

func RSimpleString(arg any) []byte {
	val := arg.(string)
	retString := fmt.Sprintf("+%s\r\n", val)
	return []byte(retString)
}

func RBulkString(arg any) []byte {
	val := arg.(string)
	valSize := len(val)
	output := fmt.Sprintf("$%d\r\n%s\r\n", valSize, val)
	return []byte(output)
}

func RNullBulkString(arg any) []byte {
	return []byte("$-1\r\n")
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
	elems := arg.([]string)
	size := len(elems)
	outputStr := fmt.Sprintf("*%d\r\n", size)
	output := []byte(outputStr)
	for _, elem := range elems {
		output = append(output, RBulkString(elem)...)
	}
	return []byte(output)
}
