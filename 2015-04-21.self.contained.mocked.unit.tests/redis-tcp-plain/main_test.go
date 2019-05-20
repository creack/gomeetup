package main

import (
	"fmt"
	"net"
	"testing"

	"github.com/garyburd/redigo/redis"
)

func newTestServer(t *testing.T, hdlr func(c net.Conn)) net.Listener {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to listen on a port: %s", err)
	}

	go func() {
		for {
			c, err := l.Accept()
			if err != nil {
				return
			}
			go hdlr(c)
		}
	}()
	return l
}

func TestRedisGet(t *testing.T) {
	ch := make(chan struct{})
	defer func() { <-ch }()

	ts := newTestServer(t, func(c net.Conn) {
		defer close(ch)
		defer c.Close()

		buf := make([]byte, 1024)
		n, err := c.Read(buf)
		if err != nil {
			t.Fatal(err)
		}
		if expect, got := "*2\r\n$3\r\nGET\r\n$3\r\nkey\r\n", string(buf[:n]); expect != got {
			t.Fatalf("Unexpected output.\nExpect:\t%q\nGot:\t%q", expect, got)
		}
		fmt.Fprintf(c, "$5\r\nvalue\r\n")
	})

	conn, err := redis.Dial("tcp", ts.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	ret, err := conn.Do("GET", "key")
	if err != nil {
		t.Fatal(err)
	}
	if expect, got := "value", string(ret.([]byte)); expect != got {
		t.Fatalf("Unexpected output.\nExpect:\t%q\nGot:\t%q", expect, got)
	}
}
