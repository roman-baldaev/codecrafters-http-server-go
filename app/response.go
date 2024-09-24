package main

import "fmt"

type Response struct {
	Protocol    string
	Status      int
	Description string
	Headers     map[string]string
	Body        string
}

func NewResponse(proto string, st int, desc string, headers map[string]string, body string) *Response {
	return &Response{
		Protocol:    proto,
		Status:      st,
		Description: desc,
		Headers:     headers,
		Body:        body,
	}
}

func (r *Response) String() string {
	if r == nil {
		return ""
	}
	headersStr := ""
	for header, value := range r.Headers {
		headersStr += fmt.Sprintf("%s: %s\r\n", header, value)
	}
	return fmt.Sprintf("%s %d %s\r\n%s\r\n%s\r\n", r.Protocol, r.Status, r.Description, headersStr, r.Body)

}
