package main

import (
	"net"
	"testing"

	dogstatsd "github.com/creack/go-dogstatsd"
)

func TestInitStatsd(t *testing.T) {
	addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	stats, err := dogstatsd.New(conn.LocalAddr().String())
	if err != nil {
		t.Fatal(err)
	}
	// Make sure the connection is established and that we receive the expected data.
	// Do not test the whole statsd client has it is already tested in the external library.
	if err := stats.Count("testcount", 42, nil, 1); err != nil {
		t.Fatal(err)
	}

	buf := make([]byte, 2048)
	n, err := conn.Read(buf)
	if err != nil {
		t.Fatal(err)
	}
	if expected, got := "testcount:42|c", string(buf[:n]); expected != got {
		t.Fatalf("Unexpected Statsd data.\nGot:\t%s\nExpect:\t%s\n", got, expected)
	}
}
