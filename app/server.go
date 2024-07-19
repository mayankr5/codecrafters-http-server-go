package main

import (
	"fmt"
	// Uncomment this block to pass the first stage
	"net"
	"os"
	"strings"
	"path"
	"bytes"
	"compress/gzip"
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

	defer conn.Close()
	buf := make([]byte, 1024)

	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error handling request: ", err.Error())
		os.Exit(1)
	}

	r := strings.Split(string(buf), "\r\n")
	m := strings.Split(r[0], " ")[0]
	p := strings.Split(r[0], " ")[1]

	var ua string
	var em string
	var response string
	var bd string
	ce := "";

	for _, val := range r {
		if strings.HasPrefix(val, "User-Agent:") {
			ua = strings.Trim(strings.Split(val, " ")[1], "\r\n")
		}
		if strings.HasPrefix(val, "Accept-Encoding:") {
			em = strings.Trim(strings.Split(val, " ")[1], "\r\n")
		}
	}


	var b bytes.Buffer
	w := gzip.NewWriter(&b)

	if m == "GET" && p == "/" {
		response = "HTTP/1.1 200 OK\r\n\r\n"
	} else if m == "GET" && strings.HasPrefix(p, "/echo") {
		bd = p[6:]
		if em == "gzip"{
			w.Write([]byte(bd))
			bd = b.String();
			ce = "Content-Encoding: gzip\r\n"
		}
		response = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n%sContent-Length:%d\r\n\r\n%s", ce, len(bd), bd)
	} else if m == "GET" && strings.HasPrefix(p, "/user-agent") {
		bd = ua
		response = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length:%d\r\n\r\n%s", len(ua), ua)
	} else if m == "GET" && strings.HasPrefix(p, "/files") {
		dir := os.Args[2]
		content, err := os.ReadFile(path.Join(dir, p[7:]))
		bd = string(content)
		if err != nil {
			response = "HTTP/1.1 404 Not Found\r\n\r\n"
		} else {
			response = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: %d\r\n\r\n%s", len(bd), string(bd))
		}
	} else if m == "POST" && strings.HasPrefix(p, "/files") {
		content := strings.Trim(r[len(r)-1], "\x00")
		dir := os.Args[2]
		_ = os.WriteFile(path.Join(dir, p[7:]), []byte(content), 0644)
		response = "HTTP/1.1 201 Created\r\n\r\n"

	} else {
		response = "HTTP/1.1 404 Not Found\r\n\r\n"
	}

	conn.Write([]byte(response))

}
