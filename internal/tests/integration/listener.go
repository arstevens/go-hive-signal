package integration_tests

import (
	"fmt"
	"math/rand"
	"net"

	"github.com/arstevens/go-request/handle"
)

type TestListener struct {
	stream chan handle.Conn
}

func newTestListener(size int) *TestListener {
	return &TestListener{
		stream: make(chan handle.Conn, size),
	}
}

func (tl *TestListener) AddConn(c handle.Conn) {
	tl.stream <- c
}

func (tl *TestListener) Accept() (handle.Conn, error) {
	conn, ok := <-tl.stream
	if !ok {
		return nil, fmt.Errorf("Channel closed in TestListener.Accept()")
	}
	return conn, nil
}

func (tl *TestListener) Close() error {
	close(tl.stream)
	return nil
}

func fakeReadRequest(conn handle.Conn) ([]byte, error) {
	c := conn.(*FakeConn)
	return c.initialData, nil
}

type FakeConn struct {
	initialData []byte
	addr        string
	closed      bool
}

func (fc *FakeConn) GetIP() net.IP { return net.ParseIP(fc.addr) }
func (fc *FakeConn) Read(b []byte) (int, error) {
	b[len(b)-1] = byte(rand.Intn(20) + 1)
	return len(b), nil
}
func (fc *FakeConn) Write([]byte) (int, error) { return 0, nil }
func (fc *FakeConn) Close() error              { fc.closed = true; return nil }
func (fc *FakeConn) IsClosed() bool            { return fc.closed }
func (fc *FakeConn) GetAddress() string        { return fc.addr }
