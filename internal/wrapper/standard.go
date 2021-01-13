package wrapper

import (
	"fmt"
	"io"
	"net"
	"time"

	"github.com/arstevens/go-request/handle"
)

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
	conn := GatewayNetConnWrapper{conn: nConn, closeCalled: false, closeDetected: false}
	return &conn, err
}

// Close closes the underlying net.Listener
func (nl *NetListener) Close() error {
	return nl.listener.Close()
}

/*GatewayNetConnWrapper implements manager.Conn and gateway.Conn,
wrapping around a standard net.Conn*/
type GatewayNetConnWrapper struct {
	conn          net.Conn
	closeCalled   bool
	closeDetected bool
}

//Read delegates to net.Conn.Read
func (gw *GatewayNetConnWrapper) Read(b []byte) (int, error) {
	return gw.conn.Read(b)
}

//Write delegates to net.Conn.Write
func (gw *GatewayNetConnWrapper) Write(b []byte) (int, error) {
	return gw.conn.Write(b)
}

//Close delegates to net.Conn.Close and records the closed state of the net.Conn
func (gw *GatewayNetConnWrapper) Close() error {
	if gw.closeCalled || gw.closeDetected {
		return fmt.Errorf("Connection already closed in GatewayNetConnWrapper.Close()")
	}
	gw.closeCalled = true
	return gw.conn.Close()
}

//GetAddress returns net.Conn.RemoteAddr().String()
func (gw *GatewayNetConnWrapper) GetAddress() string {
	addr := gw.conn.RemoteAddr()
	return addr.String()
}

//IsClosed tests whether or not the connection was closed
func (gw *GatewayNetConnWrapper) IsClosed() bool {
	if gw.closeCalled || gw.closeDetected {
		return true
	}

	one := make([]byte, 1)
	gw.conn.SetReadDeadline(time.Now())
	if _, err := gw.conn.Read(one); err == io.EOF {
		gw.closeDetected = true
	} else {
		gw.conn.SetReadDeadline(time.Time{})
	}
	return gw.closeDetected
}
