package main

import (
	"bufio"
	"io"
	"net"
	"strings"
)

type Request struct {
	// TODO: replace net.Conn by Reader interface
	Conn        net.Conn
	requestLine string
	headers     map[string]string
	body        []byte
}

func NewRequest(conn net.Conn) *Request {
	return &Request{
		Conn:    conn,
		headers: make(map[string]string),
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
	// TODO: make header in form of map[string][]string
	var splittedHeader []string
	for _, h := range headers {
		splittedHeader = strings.SplitN(h, ": ", 2)
		if len(splittedHeader) != 2 {
			continue
		}
		r.headers[splittedHeader[0]] = strings.TrimRight(splittedHeader[1], "\r\n")
	}
	return nil
}
