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
	splitStrings := strings.Split(string(buf[:size]), `\r\n`)
	count := 0
	for _, v := range splitStrings {
		if v == "PING" {
			count++
		}
	}
	return count
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage
	//
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	buf := make([]byte, 512)
	for {
		socket, err := l.Accept()
		if err != nil {
			fmt.Println("Error trying to accept tcp conn", err)
			os.Exit(1)
		}
		size, err := socket.Read(buf)
		if err != nil {
			fmt.Println("Error reading from active connection ", err)
			fmt.Printf("Closing accepted connection to %+v due to error\n", socket.RemoteAddr())
		}
		fmt.Println("Received ", (buf[:size]))
		numberOfPings := numPings(buf, size)
		if err != nil {
			fmt.Println("Error reading from active connection ", err)
			fmt.Printf("Closing accepted connection to %+v due to error\n", socket.RemoteAddr())
			socket.Close()
		}
		response := staticResponse()
		fmt.Println("static response ", string(response))
		for i := 0; i < numberOfPings; i++ {
			size, err = socket.Write(response)
			if err != nil {
				fmt.Println("Error reading from active connection ", err)
				fmt.Printf("Closing accepted connection to %+v due to error\n", socket.RemoteAddr())
				socket.Close()
			}
			fmt.Println("Wrote to socket N bytes ", size)
		}
		fmt.Println("Wrote to socket N bytes ", size)
		socket.Close()
	}
}
