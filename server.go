package customHttp

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
)

type Server struct {
	Addr    string
	Handler Handler
}

type Handler interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
}

func ListenAndServe(addr string, h Handler) error {
	if h == nil {
		h = http.DefaultServeMux
	}
	s := &Server{Addr: addr, Handler: h}
	return s.listenAndServe()
}

func (s *Server) listenAndServe() error {
	l, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}
	for {
		nc, err := l.Accept()
		if err != nil {
			return err
		}

		hc := httpConn{nc, s.Handler}

		go hc.serve()
	}
}

type httpConn struct {
	netConn net.Conn
	handler Handler
}

func (hc *httpConn) serve() {
	defer hc.netConn.Close()
	buf := bufio.NewReader(hc.netConn)
	for {
		req, err := readRequest(buf)
		if err != nil {
			retErr := fmt.Sprintf("HTTP/1.0 400 Bad Request\r\n\r\n%s\n", err.Error())
			hc.netConn.Write([]byte(retErr))
			return
		}
		res := Response{
			Headers: make(map[string]string),
			proto:   req.Proto,
		}

		keepalive, echo := req.parseConnection()
		if echo {
			res.Headers["connection"] = req.Headers["Connection"]
		}

		hc.handler.ServeHTTP(&res, req.convertor())

		res.respond(hc.netConn)
		if !keepalive {
			return
		}
	}

}

func (r *Response) respond(n net.Conn) {
	n.Write([]byte(r.makeHeader() + r.buf.String()))
}
