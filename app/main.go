package main

import (
	"flag"

	"github.com/codecrafters-io/http-server-starter-go/app/server"
)

func main() {
	directoryFlag := flag.String("directory", "", "Directory to serve files from")
	flag.Parse()

	s := server.NewServer(*directoryFlag)
	s.Listen()
	defer s.Close()

	for {
		conn := s.Accept()
		go s.HandleConnection(conn)
	}

}
