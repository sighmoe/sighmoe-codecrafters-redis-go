package main

import (
	"fmt"
)

type RedisValueType int

const (
	BULKSTR RedisValueType = iota
	ARRAY
)

type RedisValue struct {
	Len       int
	ValueType RedisValueType
	Data      []byte
}

type RedisAggregate struct {
	Len       int
	ValueType RedisValueType
	Vals      []RedisValue
}

// Parse BulkStr starting from cursor position pos and return value with cursor position after value read.
func ParseBulkStr(bs []byte, pos int) (RedisValue, int) {
	if bs[pos] != '$' {
		panic(fmt.Sprintf("Expected value \n'%s'\n to be a bulk string", string(bs[pos:])))
	}

	p := pos + 1
	len := 0
	for {
		if bs[p] == '\r' {
			break
		}
		len = (len * 10) + int(bs[p]-'0')
		p = p + 1
	}

	// fmt.Printf("Expecting bulk string of length %d\n", len)
	// skip \r\n and move cursor to bulkstr payload
	p = p + 2

	// skip \r\n bytes after payload
	return RedisValue{len, BULKSTR, bs[p:(p + len)]}, (p - pos + len + 2)
}

func ParseArray(bs []byte, pos int) (RedisAggregate, int) {
	if bs[pos] != '*' {
		panic(fmt.Sprintf("Expected value \n'%s'\n to be an array", string(bs)))
	}

	p := pos + 1
	len := 0
	for {
		if bs[p] == '\r' {
			break
		}
		len = (len * 10) + int(bs[p]-'0')
		p = p + 1
	}

	// fmt.Printf("Expecting array of length %d\n", len)
	// skip \r\n and move cursor to bulkstr payload
	p = p + 2

	var vals []RedisValue
	for i := 0; i < len; i++ {
		val, read := ParseBulkStr(bs, p)
		vals = append(vals, val)
		p = p + read
	}

	return RedisAggregate{len, ARRAY, vals}, (p - pos)
}
