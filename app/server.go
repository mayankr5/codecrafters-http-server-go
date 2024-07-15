package main

import (
	"fmt"
	// Uncomment this block to pass the first stage
	"net"
	"os"
	"strings"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage
	
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go HandleConnection(conn)
	}	
	

}

func HandleConnection(conn net.Conn) {
	req := make([]byte, 1024)
	n, err := conn.Read(req)
	if err != nil {
		fmt.Println("Failed to read the request:", err)
		return
	}

	request := string(req[:n])
	path := strings.Split(request, " ")[1]

	if path == "/" {
		conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	} else if path == "/user-agent" {
		l := len(strings.Split(request, " "))
		userAgent := strings.Trim(strings.Split(request, " ")[l-1], "\r\n")
		conn.Write([]byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(userAgent), userAgent)))
		return
	} else if strings.Split(path, "/")[1] == "files" {
		filename := strings.Trim(strings.Split(path, "/")[2], "\r\n")
		fmt.Println(filename)
		absPath := "/tmp/" + filename
		fmt.Println(absPath)
		if err == nil {
			data, _ := os.ReadFile(absPath)
			conn.Write([]byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length:%d\r\n\r\n%s", len(data), string(data))))
			return
		}
	} else if strings.Split(path, "/")[1] == "echo" {
		message := strings.Split(path, "/")[2]
		conn.Write([]byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(message), message)))
		return
	}

	conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	return
}
