package main

import (
	"log"
	"reflect"
	"strings"
	"testing"
	"time"
)

// * 1 CR LF $ 4 CR LF P I N G
var singlePing = []byte{42, 49, 13, 10, 36, 52, 13, 10, 80, 73, 78, 71, 13, 10}

// * 1 CR LF $ 4 CR LF P I N G CR LF P I N G
var doublePing = []byte{42, 49, 13, 10, 36, 52, 13, 10, 80, 73, 78, 71, 13, 10, 80, 73, 78, 71, 13, 10}

var echoCommand = []byte("*2\r\n$4\r\nECHO\r\n$3\r\nhey\r\n")

var setCmd = []byte("*3\r\n$4\r\nSET\r\n$3\r\nKEY\r\n$3\r\nVAL\r\n")

var setWithExpCmd = []byte("*5\r\n$3\r\nSET\r\n$4\r\npear\r\n$5\r\ngrape\r\n$2\r\npx\r\n$3\r\n100\r\n")

var cmdConfigGet = []byte("*3\r\n$6\r\nCONFIG\r\n$3\r\nGET\r\n$3\r\ndir\r\n")

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

func TestValidateCommandTypeSet(t *testing.T) {
	output := parseCommand(setCmd, len(setCmd))
	expected := Command{
		SET,
		3,
		[]string{"$3", "KEY", "$3", "VAL", ""},
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

func TestValidateCommandTypeConfigGet(t *testing.T) {
	output := parseCommand(cmdConfigGet, len(cmdConfigGet))
	expected := Command{
		CONFIG,
		3,
		[]string{"$3", "GET", "$3", "dir", ""},
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
func TestValidateCommandTypeSetWithExp(t *testing.T) {
	output := parseCommand(setWithExpCmd, len(setWithExpCmd))
	expected := Command{
		SET,
		5,
		[]string{"$4", "pear", "$5", "grape", "$2", "px", "$3", "100", ""},
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

func TestKeyExpirationNotExpired(t *testing.T) {
	kv := NewKvStore()
	amount := 100
	currTime := time.Now()
	futureTime := currTime.Add(time.Millisecond * time.Duration(amount))
	value := NewValue("bar", &futureTime)
	kv.SET("foo", value)
	if value.timer != futureTime {
		log.Fatal("Time is different from the value desired")
	}
	output, _ := kv.GET("foo")
	if output != "bar" {
		log.Fatalf("Expecting %s got %s", "bar", output)
	}
}
func TestKeyExpirationExpired(t *testing.T) {
	kv := NewKvStore()
	amount := 100
	currTime := time.Now()
	futureTime := currTime.Add(time.Millisecond * time.Duration(amount))
	value := NewValue("bar", &futureTime)
	kv.SET("foo", value)
	if value.timer != futureTime {
		log.Fatal("Time is different from the value desired")
	}
	time.Sleep(time.Millisecond * 101)
	output, _ := kv.GET("foo")
	if output != "" {
		log.Fatalf("Expecting %s got %s", "", output)
	}
}

func TestSerializeArrayResponse(t *testing.T) {
	respArr := []string{"dir", "/tmp/redis-files"}
	output := serializeResponse(respArr)
	expected := "*2\r\n$3\r\ndir\r\n$16\r\n/tmp/redis-files\r\n"
	if output != expected {
		log.Fatalf("Serialize Array Response expected [%s] got [%s]", expected, output)
	}
}
func TestSerializeEmptyResponse(t *testing.T) {
	respArr := []string{""}
	output := serializeResponse(respArr)
	expected := "*0\r\n"
	if output != expected {
		log.Fatalf("Serialize Array Response expected [%s] got [%s]", expected, output)
	}
}
