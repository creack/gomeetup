package main

import (
	"net"
	"testing"
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
				// TODO: handle error properly.
				return
			}
			if hdlr != nil {
				go hdlr(c)
			}
		}
	}()
	return l
}
