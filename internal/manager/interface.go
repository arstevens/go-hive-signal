package manager

import "io"

/*SwarmGateway returns a connection object to an endpoint
in a designated swarm*/
type SwarmGateway interface {
	AddEndpoint(Conn) error
	RetireEndpoint(Conn) error
	GetEndpoint() (Conn, error)
	PushEndpointAddr(string) error
	DropEndpointAddr(string) error
	GetTotalEndpoints() int
	GetEndpointAddrs() []string
	io.Closer
}

type SwarmSizeTracker interface {
	SetSize(string, int)
}

/*AgentNegotiator takes in two connection objects(the offerer
  and the acceptor) and passes session descriptions between them
  until the acceptor accepts the session*/
type AgentNegotiator func(Conn, Conn) error

//Conn represents a connection to an endpoint
type Conn interface {
	io.ReadWriteCloser
}
