package main

import (
	"bufio"
	"io"
	"net"
)

type Request struct {
	Conn        net.Conn
	requestLine string
	headers     []string
	body        []byte
}

func NewRequest(conn net.Conn) *Request {
	return &Request{
		Conn: conn,
	}
}

// Super simple HTTP request parser
func (r *Request) Parse() error {
	var err error
	reader := bufio.NewReader(r.Conn)
	r.requestLine, err = reader.ReadString('\n')

	if err != nil {
		return err
	}
	var headers []string
	var header string
	var body []byte
	totalReads := 0
	for {
		header, err = reader.ReadString('\n')
		if err != nil {
			if err == io.EOF && totalReads >= 2 {
				body = []byte(header)
				break
			}
			return err
		}
		totalReads++
		if header == "\r\n" {
			if totalReads >= 2 {
				break
			}
		}
		headers = append(headers, header)
	}
	r.body = body
	r.headers = headers
	return nil
}
