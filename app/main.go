package main

import (
	"github.com/codecrafters-io/http-server-starter-go/app/server"
)

func main() {

	s := server.NewServer()
	s.Listen()
	defer s.Close()

	for {
		conn := s.Accept()
		go s.HandleConnection(conn)
	}

}
