package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

// Ensures gofmt doesn't remove the "net" and "os" imports in stage 1 (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

func staticResponse() []byte {
	bytes := []byte{}
	msg := "+PONG\r\n"
	bytes = append(bytes, msg...)
	return bytes
}

func numPings(buf []byte, size int) int {
	splitStrings := strings.Split(string(buf[:size]), "\r\n")
	count := 0
	for _, v := range splitStrings {
		if v == "PING" {
			count++
		}
	}
	return count
}

func handleResponse(socket net.Conn) {
	buf := make([]byte, 512)
	for {
		size, err := socket.Read(buf)
		if err != nil {
			fmt.Println("Error reading from active connection ", err)
			fmt.Printf("Closing accepted connection to %+v due to error\n", socket.RemoteAddr())
			break
		}
		fmt.Printf("Received %d bytes %+v\n", size, (buf[:size]))
		// _ := numPings(buf, size)
		if err != nil {
			fmt.Println("Error reading from active connection ", err)
			fmt.Printf("Closing accepted connection to %+v due to error\n", socket.RemoteAddr())
			break
		}
		response := staticResponse()
		fmt.Println("static response ", string(response))
		size, err = socket.Write(response)
		if err != nil {
			fmt.Println("Encountered error writing to socket ", err)
			break
		}
		fmt.Println("Wrote to socket N bytes ", size)
	}
}

func main() {
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	for {
		socket, err := l.Accept()
		if err != nil {
			fmt.Println("Error trying to accept tcp conn", err)
			os.Exit(1)
		}
		go handleResponse(socket)
	}
}
