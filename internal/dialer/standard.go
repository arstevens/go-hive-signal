package dialer

import (
	"fmt"
	"io"
	"net"
	"time"

	"github.com/arstevens/go-hive-signal/internal/gateway"
)

//Implements gateway.DialEndpoint
func TCPDialEndpoint(addr string) (gateway.Conn, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	return &GatewayNetConnWrapper{
		conn:          conn,
		closeCalled:   false,
		closeDetected: false,
	}, nil
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
