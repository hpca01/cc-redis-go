package main

import (
	"log"
	"reflect"
	"strings"
	"testing"
)

// * 1 CR LF $ 4 CR LF P I N G
var singlePing = []byte{42, 49, 13, 10, 36, 52, 13, 10, 80, 73, 78, 71, 13, 10}

// * 1 CR LF $ 4 CR LF P I N G CR LF P I N G
var doublePing = []byte{42, 49, 13, 10, 36, 52, 13, 10, 80, 73, 78, 71, 13, 10, 80, 73, 78, 71, 13, 10}

var echoCommand = []byte("*2\r\n$4\r\nECHO\r\n$3\r\nhey\r\n")

func TestSinglePing(t *testing.T) {
	count := numPings(singlePing, len(singlePing))
	if count != 1 {
		log.Fatal("Error, num pings counting func is not working")
	}
}
func TestDoublePing(t *testing.T) {
	count := numPings(doublePing, len(doublePing))
	if count != 2 {
		log.Fatal("Error, num pings counting func is not working")
	}
}

func TestValidateCommandTypePing(t *testing.T) {
	output := parseCommand(singlePing, len(singlePing))
	expected := Command{
		PING,
		1,
		[]string{},
	}
	if output.command != expected.command {
		log.Fatalf("Expected type of command %+v vs %+v\n", expected.command, output.command)
	}
	if output.args != expected.args {
		log.Fatalf("Expected count of args %+v vs %+v\n", expected.args, output.args)
	}
	if reflect.DeepEqual(output.argBytes, expected.argBytes) != true {
		log.Fatalf("Expected remaining strings %+v vs %+v\n", expected.argBytes, output.argBytes)
	}
}

func TestValidateCommandTypeECHO(t *testing.T) {
	output := parseCommand(echoCommand, len(echoCommand))
	expected := Command{
		ECHO,
		2,
		[]string{"$3", "hey", ""},
	}
	if output.command != expected.command {
		log.Fatalf("Expected type of command %+v vs %+v\n", expected.command, output.command)
	}
	if output.args != expected.args {
		log.Fatalf("Expected count of args %+v vs %+v\n", expected.args, output.args)
	}
	if reflect.DeepEqual(output.argBytes, expected.argBytes) != true {
		log.Fatalf("Expected remaining strings %+v vs %+v\n", expected.argBytes, output.argBytes)
	}
}

func TestSerializeString(t *testing.T) {
	output := serializeString("bar")
	expected := "$3\r\nbar\r\n"
	if strings.Compare(output, expected) != 0 {
		log.Fatalf("Serialize string expected [%s] got [%s]", expected, output)
	}
}
