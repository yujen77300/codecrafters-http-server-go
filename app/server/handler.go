package server

import (
    "fmt"
    "net"
    "strings"
)

const CRLF = "\r\n"

func (s *Server) HandleConnection(conn net.Conn) {
    defer conn.Close()
    fmt.Println("Accepted connection from: ", conn.RemoteAddr())

    buf := make([]byte, 1024)
    n, err := conn.Read(buf)
    if err != nil {
        fmt.Println("Error reading connection: ", err)
        return
    }

    req := string(buf[:n])
    lines := strings.Split(req, CRLF)


    if len(lines) == 0 {
        handleBadRequest(conn)
        return
    }

    parts := strings.Split(lines[0], " ")
    if len(parts) < 2 {
        handleBadRequest(conn)
        return
    }

    path := parts[1]

    if path == "/" {
        handleRootPath(conn)
    } else {
        handleNotFound(conn)
    }
}

func handleBadRequest(conn net.Conn) {
    res := "HTTP/1.1 400 Bad Request\r\n\r\n"
    conn.Write([]byte(res))
}

func handleRootPath(conn net.Conn) {
    res := "HTTP/1.1 200 OK\r\n\r\n"
    conn.Write([]byte(res))
}

func handleNotFound(conn net.Conn) {
    res := "HTTP/1.1 404 Not Found\r\n\r\n"
    conn.Write([]byte(res))
}