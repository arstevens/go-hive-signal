package connect

import (
	"net"

	"github.com/arstevens/go-request/handle"
)

type IdentityVerifier interface {
	// IP of requester, originID, is a long on request
	Analyze(net.IP, string, bool) bool
}

type SwarmConnector interface {
	// Connection code, connection object to requester
	ProcessConnection(int, handle.Conn) error
}

type ConnectionRequest interface {
	GetRequestCode() int
	GetOriginID() string
	IsLogOn() bool
}

type NetConn interface {
	handle.Conn
	GetIP() net.IP
}
