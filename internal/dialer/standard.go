package dialer

import (
	"net"

	"github.com/arstevens/go-hive-signal/internal/manager"
)

var NewConnWraper func(net.Conn) manager.Conn = nil

//Implements gateway.DialEndpoint
func TCPDialEndpoint(addr string) (manager.Conn, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return NewConnWraper(conn), nil
}
