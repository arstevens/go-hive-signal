package route

import (
	"io"
	"net"

	"github.com/arstevens/go-request/handle"
)

/* Listener defines a type that can accept new connections
to receive requests */
type Listener interface {
	io.Closer
	Accept() (handle.Conn, error)
}

// NetListener wraps a net.Listener so it implements the Listener interface
type NetListener struct {
	listener net.Listener
}

func NewNetListener(listener net.Listener) *NetListener {
	return &NetListener{
		listener: listener,
	}
}

// Accept returns a net.Conn wrapped as a NetConn and an error
func (nl *NetListener) Accept() (handle.Conn, error) {
	nConn, err := nl.listener.Accept()
	conn := NetConn{conn: nConn}
	return &conn, err
}

// Close closes the underlying net.Listener
func (nl *NetListener) Close() error {
	return nl.listener.Close()
}

// NetConn wraps a net.Conn so it implements the Conn interface
type NetConn struct {
	conn net.Conn
}

// Read reads data from the net.Conn into a slice of bytes
func (nc *NetConn) Read(b []byte) (int, error) {
	return nc.conn.Read(b)
}

// Write writes data from a slice of bytes to the underlying net.Conn
func (nc *NetConn) Write(b []byte) (int, error) {
	return nc.conn.Write(b)
}

// Close closes the underlying net.Conn
func (nc *NetConn) Close() error {
	return nc.conn.Close()
}
