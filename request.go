package customHttp

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Request struct {
	Method  string
	URI     string
	Proto   string
	Headers map[string]string

	Body io.ReadCloser
}

func readRequest(buf *bufio.Reader) (*Request, error) {
	req := Request{
		Headers: make(map[string]string),
	}

	// Read and Parse first line of http request (request line)
	{
		requestLine, err := reqReader(buf)
		if err != nil {
			return nil, err
		}
		if err := parseRequestLine(requestLine, &req); err != nil {
			return nil, err
		}
	}
	// Read and Parse all Headers
	for {
		ln, err := reqReader(buf)
		if err != nil {
			return nil, err
		}
		if len(ln) == 0 {
			break
		}
		if err := parseHeaderLine(ln, req.Headers); err != nil {
			return nil, err
		}
	}

	return &req, nil
}

func reqReader(buf *bufio.Reader) (string, error) {
	ln, err := buf.ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.TrimSuffix(ln, "\r\n"), nil
}

func parseRequestLine(rql string, r *Request) error {
	s := strings.Split(rql, " ")
	if len(s) != 3 {
		return fmt.Errorf("malformed HTTP request, %s", s)
	}

	r.Method = s[0]
	r.URI = s[1]
	r.Proto = s[2]

	return nil
}

func parseHeaderLine(rql string, headerMap map[string]string) error {
	s := strings.SplitN(rql, ":", 2)
	if len(s) != 2 {
		return fmt.Errorf("malformed HTTP header line: %s", rql)
	}
	key := strings.TrimSpace(s[0])
	value := strings.TrimSpace(s[1])
	headerMap[key] = value
	return nil
}

func (req *Request) parseConnection() (bool, bool) {
	conn := strings.ToLower(req.Headers["Connection"])
	if req.Proto == "HTTP/1.0" && conn == "keep-alive" {
		return true, true
	}
	if req.Proto == "HTTP/1.1" && conn == "close" {
		return false, true
	}

	return false, false
}

func (r *Request) convertor() *http.Request {
	req, _ := http.NewRequest(r.Method, r.URI, r.Body)

	req.Proto = r.Proto

	for key, value := range r.Headers {
		req.Header.Set(key, value)
	}

	return req
}
