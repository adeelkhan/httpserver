package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

const (
	SERVER_HOST = "localhost"
	SERVER_PORT = "1111"
	SERVER_TYPE = "tcp"
)

func main() {
	server, err := net.Listen(SERVER_TYPE, SERVER_HOST+":"+SERVER_PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	defer server.Close()
	fmt.Println("listening on " + SERVER_HOST + ":" + SERVER_PORT)

	for {
		connection, err := server.Accept()
		if err != nil {
			fmt.Println("Error accepting:", err.Error())
			os.Exit(1)
		}
		go processHttpRequest(connection)
	}
}

func readFile(f string) ([]byte, int) {
	file, err := ioutil.ReadFile("./public/" + f)
	if err != nil {
		log.Printf("failed reading file: %s", err)
		return []byte{}, 0
	}
	return file, len(file)
}
func processHeader(conn net.Conn) string {
	buffer := make([]byte, 1024)
	mLen, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	path := strings.Split(string(buffer[:mLen]), "\n")[0]
	method, url := strings.Split(path, " ")[0], strings.Split(path, " ")[1]

	if method == "GET" {
		return url[1:]
	}
	return ""
}
func processHttpRequest(conn net.Conn) {
	defer conn.Close()
	url := processHeader(conn)
	if url == "" {
		response404Error(conn)
		return
	}
	outFile, len := readFile(url)
	if len == 0 {
		response404Error(conn)
		return
	}
	str := strconv.Itoa(len)
	response := "HTTP/1.1 200 OK\nContent-Type: text/html\nContent-Length: " + str + "\n\n" + string(outFile)
	_, err := conn.Write([]byte(string(response)))
	if err != nil {
		fmt.Println("Cant able to send response")
	}
}
func response404Error(conn net.Conn) {
	response := "HTTP/1.1 404 OK\nContent-Type: text/plain\n\n\n"
	_, err := conn.Write([]byte(string(response)))
	if err != nil {
		fmt.Println("Cant able to send response")
	}
}
