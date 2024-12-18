package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

// Ensures gofmt doesn't remove the "net" and "os" imports above (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

func main() {
	// run server
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	defer l.Close()

	// accept connection
	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
	request := NewRequest(conn)
	err = request.Parse()
	if err != nil {
		fmt.Printf("Error reading request = %s\n", err.Error())
		os.Exit(1)
	}
	fmt.Println(request.headers)
	fmt.Println(request.requestLine)
	requestData := strings.Split(request.requestLine, " ")
	var resp *Response
	if len(requestData) < 2 {
		resp = NewResponse("HTTP/1.1", 404, "Not Found", nil, "")
		_, err = conn.Write([]byte(resp.String()))
		if err != nil {
			fmt.Println("Error writing response: ", err.Error())
			os.Exit(1)
		}
		return
	}
	path := strings.Split(requestData[1], "/")
	switch len(path) {
	case 0:
		fmt.Println(path)
		resp = NewResponse("HTTP/1.1", 200, "OK", nil, "")
	case 1:
		fmt.Println(path)
	case 2:
		if path[0] == path[1] && path[0] == "" {
			resp = NewResponse("HTTP/1.1", 200, "OK", nil, "")
			break
		}
		if path[1] == "user-agent" {
			val := request.headers["User-Agent"]
			resp = NewResponse("HTTP/1.1", 200, "OK", map[string]string{"Content-Type": "text/plain", "Content-Length": fmt.Sprintf("%d", len([]byte(val)))}, val)
			break
		}
		resp = NewResponse("HTTP/1.1", 404, "Not Found", nil, "")
	case 3:
		if path[1] == "echo" {
			resp = NewResponse("HTTP/1.1", 200, "OK", map[string]string{"Content-Type": "text/plain", "Content-Length": fmt.Sprintf("%d", len([]byte(path[2])))}, path[2])
		} else {
			resp = NewResponse("HTTP/1.1", 404, "Not Found", nil, "")
		}
	default:
		resp = NewResponse("HTTP/1.1", 404, "Not Found", nil, "")
	}

	_, err = conn.Write([]byte(resp.String()))
	defer conn.Close()
	if err != nil {
		fmt.Println("Error writing response: ", err.Error())
		os.Exit(1)
	}
}
