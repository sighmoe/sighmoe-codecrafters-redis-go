package main

import (
	"fmt"
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestParseBulkStr(t *testing.T) {
	test := "$3\r\nhey\r\n"
	expected := RedisValue{3, BULKSTR, []byte{'h', 'e', 'y'}}
	actual, pos := ParseBulkStr([]byte(test), 0)
	if !cmp.Equal(expected, actual) {
		t.Fatalf(`Expected %v but got %v`, expected, actual)
	}

	if pos != len(test) {
		t.Fatalf(`Expected 9 bytes to be read but got %v`, pos)
	}
}

func TestParseBulkStrLong(t *testing.T) {
	str := "heythisisamuchlongerstring4testing"
	test := fmt.Sprintf("$%d\r\n%s\r\n", len(str), str)
	expected := RedisValue{len(test), BULKSTR, []byte(test)}
	actual, pos := ParseBulkStr([]byte(test), 0)
	if !cmp.Equal(expected, actual) {
		t.Fatalf(`Expected %v but got %v`, expected, actual)
	}

	if pos != len(test) {
		t.Fatalf(`Expected %d bytes to be read but got %v`, 7+len(test), pos)
	}
}

func TestParseArray(t *testing.T) {
	test := "*2\r\n$4\r\necho\r\n$3\r\nhey\r\n"
	val1 := RedisValue{4, BULKSTR, []byte{'e', 'c', 'h', 'o'}}
	val2 := RedisValue{3, BULKSTR, []byte{'h', 'e', 'y'}}
	expected := RedisAggregate{len(test), ARRAY, []RedisValue{val1, val2}}
	actual, pos := ParseArray([]byte(test), 0)
	if !cmp.Equal(expected, actual) {
		t.Fatalf(`Expected %v but got %v`, expected, actual)
	}

	// add constant 7 for $xx (3 bytes) and two sets of \r\n bytes
	if pos != len(test) {
		t.Fatalf(`Expected %d bytes to be read but got %v`, len(test), pos)
	}
}
