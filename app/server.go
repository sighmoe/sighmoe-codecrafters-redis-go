package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go HandleConnection(conn)
	}
}

func HandleConnection(conn net.Conn) {
	defer conn.Close()
	for {
		buf := make([]byte, 1024)
		//fmt.Printf("\nReading from connection...\n")
		_, err := conn.Read(buf)
		if err != nil {
			fmt.Printf("\nGot err reading from connection: %v. Terminating connection...\n", err.Error())
			break
		}
		//fmt.Printf("Read\n\n%s\nfrom the connection\n", string(buf))
		value, _ := ParseArray(buf, 0)
		//fmt.Printf("\nGot value:\n%#v", value)

		request := RequestFromValue(value)
		HandleRequest(request, conn)
	}
}

type RedisRequest struct {
	request string
	command string
	args    []string
}

func RequestFromValue(v RedisAggregate) RedisRequest {
	var strs []string

	for _, val := range v.Vals {
		strs = append(strs, strings.ToLower(string(val.Data)))
	}

	request := strings.Join(strs, " ")
	command := strs[0]
	return RedisRequest{request, command, strs[1:]}
}

func HandleRequest(r RedisRequest, conn net.Conn) {
	switch r.command {
	case "ping":
		conn.Write([]byte("+PONG\r\n"))
	case "echo":
		str := strings.Join(r.args, "")
		response := fmt.Sprintf("$%v\r\n%s\r\n", len(str), str)
		conn.Write([]byte(response))
	default:
		panic(fmt.Sprintf("Received unknown command: %s", r.command))

	}

}
