package server

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
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

	switch reqStatusLine[0] {
	case "GET":
		if path == "/" {
			handleRootPath(conn)
		} else if strings.HasPrefix(path, "/echo/") {
			echoStr := strings.TrimPrefix(path, "/echo/")
			handleEchoPath(conn, echoStr)
		} else if path == "/user-agent" {
			handleUserAgent(conn, lines)
		} else if strings.HasPrefix(path, "/files/") {
			fileName := strings.TrimPrefix(path, "/files/")
			s.handleFilesPath(conn, fileName)
		} else {
			handleNotFound(conn)
		}
	case "POST":
		if strings.HasPrefix(path, "/files/") {
			// Find the position where headers end and body begins (marked by two consecutive CRLFs)
			bodyStart := strings.Index(req, CRLF+CRLF) + 4

			requestBody := ""
			// Check if there's any content after headers
			if bodyStart < len(req) {
				requestBody = req[bodyStart:]
			}
			fmt.Println("Request body: ", requestBody)
			fileName := strings.TrimPrefix(path, "/files/")
			fmt.Printf("POST request for file: %s with body length: %d\n", fileName, len(requestBody))
			s.handlePostFilePath(conn, fileName, requestBody)
		} else {
			handleNotFound(conn)
		}

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

func (s *Server) handleFilesPath(conn net.Conn, fileName string) {
	filePath := filepath.Join(s.fileDirctory, fileName)

	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			handleNotFound(conn)
		} else {
			fmt.Printf("Error checking file: %v\n", err)
			handleBadRequest(conn)
		}
		return
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		handleBadRequest(conn)
		return
	}

	statusLine := "HTTP/1.1 200 OK\r\n"

	contentLength := len(content)
	headers := fmt.Sprintf("Content-Type: application/octet-stream\r\nContent-Length: %d\r\n\r\n", contentLength)

	response := statusLine + headers + string(content)
	conn.Write([]byte(response))

	fmt.Printf("File response sent: %s (%d bytes)\n", fileName, contentLength)
}

func (s *Server) handlePostFilePath(conn net.Conn, fileName string, requestBody string) {
	filePath := filepath.Join(s.fileDirctory, fileName)

	err := os.WriteFile(filePath, []byte(requestBody), 0644)
	if err != nil {
		fmt.Printf("Error writing to file: %v\n", err)
		handleBadRequest(conn)
		return
	}

	statusLine := "HTTP/1.1 201 Created\r\n\r\n"
	conn.Write([]byte(statusLine))

	fmt.Printf("File created: %s (%d bytes)\n", fileName, len(requestBody))
}
