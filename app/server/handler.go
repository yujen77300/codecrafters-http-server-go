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

	reqStatusLine := strings.Split(lines[0], " ")
	if len(reqStatusLine) < 2 {
		handleBadRequest(conn)
		return
	}
	fmt.Println("Request line: ", reqStatusLine[0], reqStatusLine[1], reqStatusLine[2])

	path := reqStatusLine[1]

	fmt.Println("Request path: ", path)

	if path == "/" {
		handleRootPath(conn)
	} else if strings.HasPrefix(path, "/echo/") {
		echoStr := strings.TrimPrefix(path, "/echo/")
		handleEchoPath(conn, echoStr)
	} else if path == "/user-agent"{
		handleUserAgent(conn, lines)
	}else{
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


func handleUserAgent(conn net.Conn, lines []string) {
    userAgent := ""
    for _, line := range lines {
        if strings.HasPrefix(line, "User-Agent:") {
            userAgent = strings.TrimPrefix(line, "User-Agent: ")
            break
        }
    }

    if userAgent == "" {
        handleBadRequest(conn)
        return
    }

    statusLine := "HTTP/1.1 200 OK\r\n"

    contentLength := len(userAgent)
    headers := fmt.Sprintf("Content-Type: text/plain\r\nContent-Length: %d\r\n\r\n", contentLength)

    body := userAgent

    response := statusLine + headers + body
    conn.Write([]byte(response))

    fmt.Printf("User-Agent response sent: %s\n", userAgent)
}