package connector

import (
	"net"

	"github.com/arstevens/go-request/handle"
)

/*IdentityVerifier takes in information about a request and
requester and returns whether or not the request is valid*/
type IdentityVerifier interface {
	// IP of requester, originID, is a log on request
	Analyze(net.IP, string, bool) bool
}

//SwarmConnector connects a connection to a swarm
type SwarmConnector interface {
	// SwarmID if exists, Connection code, connection object to requester
	ProcessConnection(string, bool, handle.Conn) error
}

/*ConnectionRequest is the request type that a
ConnectionHandler can process*/
type ConnectionRequest interface {
	GetSwarmID() string
	GetOriginID() string
	IsLogOn() bool
}

/*NetConn is a type of handle.Conn that has an additional
method GetIP() since it represents a network connection*/
type NetConn interface {
	handle.Conn
	GetIP() net.IP
}
