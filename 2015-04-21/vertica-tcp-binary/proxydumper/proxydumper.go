package proxydumper

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"sync"
)

// Dump generates a Go [][]byte{} variable with the data read from `r`.
func Dump(r io.Reader, w, logger io.Writer) error {
	fmt.Fprintf(logger, "var data = [][]byte{\n")
	defer func() { fmt.Fprintf(logger, "}\n") }()
	for {
		buf := make([]byte, 32*4096)
		n, err := r.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		tmp := bytes.NewBuffer([]byte("        {%s%s"))
		splited := false
		for i := 0; i < n; i++ {
			if i != 0 && i%10 == 0 { // avoid long lines
				splited = true
				tmp = bytes.NewBuffer(bytes.TrimRight(tmp.Bytes(), " "))
				fmt.Fprintf(tmp, "\n\t\t")
			}
			fmt.Fprintf(tmp, "0x%02x, ", buf[i])
		}
		tmp = bytes.NewBuffer(bytes.TrimRight(tmp.Bytes(), ", "))
		fmt.Fprintf(tmp, "%%s},\n")
		if splited {
			fmt.Fprintf(logger, tmp.String(), "\n", "\t\t", ",\n\t")
		} else {
			fmt.Fprintf(logger, tmp.String(), "", "", "")
		}
		if _, err := w.Write(buf[:n]); err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
	}
	return nil
}

// ProxyListener wraps net.Listener with stopped flag and string Addr getters.
type ProxyListener struct {
	net.Listener
	Remote  net.Conn
	lock    sync.RWMutex
	stopped bool
}

// Close wraps underlying Close and set the stopped flag
func (p *ProxyListener) Close() error {
	p.lock.Lock()
	p.stopped = true
	p.lock.Unlock()
	return p.Listener.Close()
}

// HostLocal return the address which the proxy listen on locally
func (p *ProxyListener) HostLocal() string { return strings.Split(p.Listener.Addr().String(), ":")[0] }

// PortLocal return the port which the proxy listen on locally
func (p *ProxyListener) PortLocal() string { return strings.Split(p.Listener.Addr().String(), ":")[1] }

// AddrLocal return the address:port which the proxy listen on locally
func (p *ProxyListener) AddrLocal() string { return p.Listener.Addr().String() }

// Start creates a MITM logger between localAddr and remoteAddr
// NOTE: the given log paths will be overridden.
func Start(network, localAddr, remoteAddr, clientLogPath, serverLogPath string) (*ProxyListener, chan error, error) {
	remoteConn, err := net.Dial(network, remoteAddr)
	if err != nil {
		return nil, nil, err
	}
	l, err := net.Listen(network, localAddr)
	if err != nil {
		return nil, nil, err
	}
	clientLog, err := os.Create(clientLogPath)
	if err != nil {
		return nil, nil, err
	}
	serverLog, err := os.Create(serverLogPath)
	if err != nil {
		return nil, nil, err
	}

	p := &ProxyListener{
		Remote:   remoteConn,
		Listener: l,
	}
	errChan := make(chan error, 3)
	go func() {
		// TOOD: this is risky, might endup send on closed chan.
		defer close(errChan)
		for {
			c, err := l.Accept()
			if err != nil {
				p.lock.RLock()
				if p.stopped {
					return
				}
				p.lock.RUnlock()
				errChan <- err
				return
			}
			go func() {
				ch := make(chan struct{}, 2)
				defer close(ch)
				go func() {
					defer func() { ch <- struct{}{} }()
					errChan <- Dump(c, remoteConn, clientLog)
				}()
				go func() {
					defer func() { ch <- struct{}{} }()
					errChan <- Dump(remoteConn, c, serverLog)
				}()
				<-ch
				<-ch
			}()
		}
	}()
	return p, errChan, nil
}
