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
	fmt.Println("Request line: ", parts[0], parts[1], parts[2])

	path := parts[1]

	fmt.Println("Request path: ", path)

	if path == "/" {
		handleRootPath(conn)
	} else if strings.HasPrefix(path, "/echo/") {
		echoStr := strings.TrimPrefix(path, "/echo/")
		handleEchoPath(conn, echoStr)
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

func handleEchoPath(conn net.Conn, content string) {
    // Status line
    statusLine := "HTTP/1.1 200 OK\r\n"

    // Headers
    contentLength := len(content)
    headers := fmt.Sprintf("Content-Type: text/plain\r\nContent-Length: %d\r\n\r\n", contentLength)

    // Response body
    body := content

    response := statusLine + headers + body
    conn.Write([]byte(response))

    fmt.Printf("Echo response sent: %s\n", content)
}