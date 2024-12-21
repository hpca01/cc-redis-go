package main

import (
	"log"
	"testing"
)

// * 1 CR LF $ 4 CR LF P I N G
var singlePing = []byte{42, 49, 13, 10, 36, 52, 13, 10, 80, 73, 78, 71, 13, 10}

var doublePing = []byte{42, 49, 13, 10, 36, 52, 13, 10, 80, 73, 78, 71, 13, 10, 80, 73, 78, 71, 13, 10}

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

func TestStaticResponse() {

}
