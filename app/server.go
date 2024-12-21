package main

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

// Ensures gofmt doesn't remove the "net" and "os" imports in stage 1 (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

type CommandType int32

const (
	PING CommandType = iota
	ECHO
	SET
	GET
)

type Command struct {
	command  CommandType
	args     int
	argBytes []string
}

type RdbArgs struct {
	dir        string
	dbfilename string
}

type ListStr []string

func (s *ListStr) has(search string) bool {
	size := len(*s)
	for i := 0; i < size; i++ {
		if (*s)[i] == search {
			return true
		}
	}
	return false
}
func (s *ListStr) idxOf(search string) int {
	size := len(*s)
	for i := 0; i < size; i++ {
		if (*s)[i] == search {
			return i
		}
	}
	return -1
}

func NewRdbArgsFromCmdArgs(args []string) *RdbArgs {
	var argsList ListStr = args
	var dirName string
	var dbFileName string
	if argsList.has("--dir") && argsList.has("--dbfilename") {
		dir := argsList.idxOf("--dir")
		dirName = args[dir+1]
		dbfile := argsList.idxOf("--dbfilename")
		dbFileName = args[dbfile+1]
	}
	return &RdbArgs{dirName, dbFileName}
}

const ok = "+OK\r\n"
const notExist = "$-1\r\n"
const emptyResponse = "*0\r\n"

var KvStore KeyValueStore

func init() {
	KvStore = *NewKvStore()
}

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

func parseCommand(buf []byte, size int) Command {
	asStr := string(buf[:size])
	splitStrings := strings.Split(asStr, "\r\n")
	// first part is always the num of args
	numArgs, err := strconv.Atoi(splitStrings[0][1:])
	if err != nil {
		fmt.Println("Error convering num args")
		panic(err)
	}
	var command Command
	switch strings.ToUpper(splitStrings[2]) {
	case "PING":
		command = Command{PING, numArgs, []string{}}
	case "ECHO":
		// ECHO always has something in args
		command = Command{ECHO, numArgs, splitStrings[3:]}
	case "SET":
		command = Command{SET, numArgs, splitStrings[3:]}
	case "GET":
		command = Command{GET, numArgs, splitStrings[3:]}
	}
	return command
}

func handlePing(socket net.Conn, command Command) {
	_ = command // for future
	response := staticResponse()
	fmt.Println("static response ", string(response))
	_, err := socket.Write(response)
	if err != nil {
		fmt.Println("Encountered error writing to socket ", err)
	}
}

func handleEcho(socket net.Conn, command Command) {
	response := strings.Join(command.argBytes, "\r\n")
	fmt.Println("Echo response ", response)
	_, err := socket.Write([]byte(response))
	if err != nil {
		fmt.Println("Encountered error writing ECHO response to socket ", err)
	}
}

func serializeString(str string) string {
	len := len(str)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("$%d", len))
	sb.WriteString(fmt.Sprintf("\r\n%s\r\n", str))
	return sb.String()
}

func serializeResponse(arr []string) string {
	size := len(arr)
	if size == 0 {
		return emptyResponse
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("*%d\r\n", size))
	for i := 0; i < size; i++ {
		sb.WriteString(serializeString(arr[i]))
	}
	return sb.String()
}

func handleGet(socket net.Conn, command Command) {
	// this assumes that format is *2\r\n$4\r\GET\r\n$3\r\nKEY\r\n
	key := command.argBytes[1]
	value, err := KvStore.GET(key)
	if errors.Is(err, ErrKeyNotFound) || errors.Is(err, ErrKeyExpired) {
		_, err := socket.Write([]byte(notExist))
		if err != nil {
			fmt.Println("Encountered error writing GET response to socket ", err)
		}
		return
	}
	output := serializeString(value)
	_, err = socket.Write([]byte(output))
	if err != nil {
		fmt.Println("Encountered error writing GET response to socket ", err)
	}
}

func handleSet(socket net.Conn, command Command) {
	// this assumes that format is *3\r\n$4\r\nSET\r\n$3\r\nKEY\r\n$3\r\nVAL\r\n
	key := command.argBytes[1]
	value := command.argBytes[3]
	if command.args < 5 {
		KvStore.SET(key, NewValue(value, nil))
	} else {
		amtTime, err := strconv.Atoi(command.argBytes[len(command.argBytes)-2])
		if err != nil {
			panic("Error converting px value to time in handle set")
		}
		futureTime := time.Now().Add(time.Millisecond * time.Duration(amtTime))
		KvStore.SET(key, NewValue(value, &futureTime))
	}
	_, err := socket.Write([]byte(ok))
	if err != nil {
		fmt.Println("Encountered error writing SET response to socket ", err)
	}
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
		command := parseCommand(buf, size)
		switch command.command {
		case PING:
			handlePing(socket, command)
		case ECHO:
			handleEcho(socket, command)
		case GET:
			handleGet(socket, command)
		case SET:
			handleSet(socket, command)
		}
	}
	socket.Close()
}

func main() {
	fmt.Println("Logs from your program will appear here!")
	args := os.Args[1:]
	rdb := NewRdbArgsFromCmdArgs(args)
	fmt.Printf("Got the following command line args %+v\n", rdb)

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
