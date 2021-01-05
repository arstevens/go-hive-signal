package integration_tests

import (
	"fmt"
	"net"

	"github.com/arstevens/go-request/handle"
)

type TestListener struct {
	stream chan handle.Conn
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
	return c.initalData, nil
}

type FakeConn struct {
	initalData []byte
	addr       string
	closed     bool
}

func (fc *FakeConn) GetIP() net.IP             { return net.ParseIP(fc.addr) }
func (fc *FakeConn) Read([]byte) (int, error)  { return 0, nil }
func (fc *FakeConn) Write([]byte) (int, error) { return 0, nil }
func (fc *FakeConn) Close() error              { fc.closed = true; return nil }
func (fc *FakeConn) IsClosed() bool            { return fc.closed }
func (fc *FakeConn) GetAddress() string        { return fc.addr }
