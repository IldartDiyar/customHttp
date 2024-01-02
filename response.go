package customHttp

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type Response struct {
	Status  int
	Headers map[string]string
	proto   string
	buf     bytes.Buffer
}

func (r *Response) WriteHeader(statusCode int) {
	r.Status = statusCode
}

func (r *Response) Write(b []byte) (int, error) {
	return r.buf.Write(b)
}
func (r *Response) Header() http.Header {
	httpHeader := make(http.Header)
	for key, value := range r.Headers {
		httpHeader.Set(key, value)
	}
	return httpHeader
}

func (r *Response) makeHeader() string {
	r.Headers["Date"] = time.Now().UTC().Format("Mon, 02 Jan 2006 15:04:05 MST")
	r.Headers["Content-Length"] = strconv.Itoa(r.buf.Len())
	headers := fmt.Sprintf("%s %d %s\r\n", r.proto, r.Status, http.StatusText(r.Status))
	for k, v := range r.Headers {
		headers += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	headers += "\r\n"
	return headers
}
