package server

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const CRLF = "\r\n"

type HttpRequest struct {
	Method  string
	URL     string
	Version string
	Headers map[string]string
	Body    string
}

type HttpResponse struct {
	Status  int
	Version string
	Headers map[string]string
	Body    []byte
}

func (s *Server) HandleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Println("Accepted connection from: ", conn.RemoteAddr())

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading connection: ", err)
		return
	}

	httpReq := parseRequest(buf[:n])

	response := s.route(httpReq)

	encResponse := response.Encode()

	conn.Write(encResponse)

}

func parseRequest(bs []byte) *HttpRequest {
	req := string(bs)
	lines := strings.Split(req, CRLF)
	if len(lines) == 0 {
		return nil
	}

	reqStatusLine := strings.Split(lines[0], " ")
	if len(reqStatusLine) < 3 {
		return nil
	}

	headers := make(map[string]string)
	headerEnd := 0
	for i := 1; i < len(lines); i++ {
		if lines[i] == "" {
			headerEnd = i
			break
		}

		parts := strings.SplitN(lines[i], ": ", 2)
		if len(parts) == 2 {
			headers[parts[0]] = parts[1]
		}
	}

	body := ""
	if headerEnd+1 < len(lines) {
		body = strings.Join(lines[headerEnd+1:], CRLF)
	}

	return &HttpRequest{
		Method:  reqStatusLine[0],
		URL:     reqStatusLine[1],
		Version: reqStatusLine[2],
		Headers: headers,
		Body:    body,
	}
}

func (s *Server) route(req *HttpRequest) *HttpResponse {
	if req == nil {
		return buildResponse(400, "HTTP/1.1")
	}

	if req.URL == "/" {
		return buildResponse(200, req.Version)
	} else if strings.HasPrefix(req.URL, "/echo/") {
		echoStr := strings.TrimPrefix(req.URL, "/echo/")
		return buildResponseWithBody(200, req.Version, []byte(echoStr), "text/plain")
	} else if req.URL == "/user-agent" {
		userAgent, exists := req.Headers["User-Agent"]
		if !exists {
			return buildResponse(400, req.Version)
		}
		return buildResponseWithBody(200, req.Version, []byte(userAgent), "text/plain")
	} else if strings.HasPrefix(req.URL, "/files/") {
		fileName := strings.TrimPrefix(req.URL, "/files/")

		if req.Method == "GET" {
			return s.handleGetFile(fileName, req.Version)
		} else if req.Method == "POST" {
			return s.handlePostFile(fileName, req.Body, req.Version)
		}
	}

	return buildResponse(404, req.Version)
}

func buildResponse(status int, version string) *HttpResponse {
	return &HttpResponse{
		Status:  status,
		Version: version,
		Headers: make(map[string]string),
		Body:    []byte(""),
	}
}

func buildResponseWithBody(status int, version string, body []byte, contentType string) *HttpResponse {
	headers := make(map[string]string)
	headers["Content-Type"] = contentType
	headers["Content-Length"] = strconv.Itoa(len(body))

	return &HttpResponse{
		Status:  status,
		Version: version,
		Headers: headers,
		Body:    body,
	}
}

func (s *Server) handleGetFile(fileName string, version string) *HttpResponse {
	filePath := filepath.Join(s.fileDirctory, fileName)

	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return buildResponse(404, version)
		} else {
			fmt.Printf("Error checking file: %v\n", err)
			return buildResponse(400, version)
		}
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return buildResponse(400, version)
	}

	fmt.Printf("File response sent: %s (%d bytes)\n", fileName, len(content))
	return buildResponseWithBody(200, version, content, "application/octet-stream")
}

func (s *Server) handlePostFile(fileName string, requestBody string, version string) *HttpResponse {
	filePath := filepath.Join(s.fileDirctory, fileName)

	err := os.WriteFile(filePath, []byte(requestBody), 0644)
	if err != nil {
		fmt.Printf("Error writing to file: %v\n", err)
		return buildResponse(400, version)
	}

	fmt.Printf("File created: %s (%d bytes)\n", fileName, len(requestBody))
	return buildResponse(201, version)
}

func (r *HttpResponse) Encode() []byte {
	response := make([]byte, 0)

	statusText := getStatusText(r.Status)
	statusLine := fmt.Sprintf("%s %d %s\r\n", r.Version, r.Status, statusText)
	response = append(response, statusLine...)

	for k, v := range r.Headers {
		headerLine := fmt.Sprintf("%s: %s\r\n", k, v)
		response = append(response, headerLine...)
	}

	// headers 和 body 之間的空行
	response = append(response, "\r\n"...)

	response = append(response, r.Body...)

	return response
}

func getStatusText(status int) string {
	switch status {
	case 200:
		return "OK"
	case 201:
		return "Created"
	case 400:
		return "Bad Request"
	case 404:
		return "Not Found"
	default:
		return ""
	}
}
