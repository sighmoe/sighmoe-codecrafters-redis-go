package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

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
		fmt.Printf("\nReading from connection...\n")
		_, err := conn.Read(buf)
		if err != nil {
			fmt.Printf("\nGot err reading from connection: %v. Terminating connection...\n", err.Error())
			break
		}
		fmt.Printf("Read\n\n%s\nfrom the connection\n", string(buf))
		conn.Write([]byte("+PONG\r\n"))
	}
}
