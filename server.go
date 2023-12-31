package customHttp

import (
	"bufio"
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
}
