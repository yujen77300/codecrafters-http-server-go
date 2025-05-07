package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

const CRLF = "\r\n"

type Server struct {
	listener net.Listener
}

func main() {
	s := Server{}
	s.Start()
}

func (s *Server) Start() {
	s.Listen()
	defer s.Close()
	conn := s.Accept()
	fmt.Println("Accepted connection from: ", conn.RemoteAddr())

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error accepting connection: ", err)
	}

	req := string(buf[:n])
	lines := strings.Split(req, CRLF)
	path := strings.Split(lines[0], " ")[1]
	var res string
	if path == "/" {
		res = "HTTP/1.1 200 OK\r\n\r\n"
	} else {
		res = "HTTP/1.1 404 Not Found\r\n\r\n"
	}
	conn.Write([]byte(res))

}

func (s *Server) Listen() {
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	s.listener = l
}
func (s *Server) Accept() net.Conn {
	conn, err := s.listener.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
	return conn
}
func (s *Server) Close() {
	err := s.listener.Close()
	if err != nil {
		fmt.Println("Failed to close listener:", err.Error())
	}
}
