package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

// Ensures gofmt doesn't remove the "net" and "os" imports above (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage
	//
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	defer l.Close()

	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
	request, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("Error reading request")
		os.Exit(1)
	}
	requestData := strings.Split(request, " ")
	var output string
	if requestData[1] != "/" {
		output = "HTTP/1.1 404 Not Found\r\n\r\n"
	} else {
		output = "HTTP/1.1 200 OK\r\n\r\n"
	}
	_, err = conn.Write([]byte(output))
	defer conn.Close()
	if err != nil {
		fmt.Println("Error writing response: ", err.Error())
		os.Exit(1)
	}
}
