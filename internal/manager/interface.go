package manager

import "io"

/*SwarmGateway returns a connection object to an endpoint
in a designated swarm*/
type SwarmGateway interface {
	PushEndpoint(string) error
	RemoveEndpoint(string) error
	//Returns the connection, the preferred load parameter, and an error
	GetEndpoint() (Conn, int, error)
	GetTotalEndpoints() int
	GetEndpointAddrs() []string
	io.Closer
}

//An object that can create new SwarmGateways
type SwarmGatewayGenerator interface {
	New() SwarmGateway
}

//Defines an object that keeps track of swarm sizes
type SwarmInfoTracker interface {
	AddPreferredLoadDatapoint(string, int)
	SetSize(string, int)
	Delete(string)
}

/*AgentNegotiator takes in two connection objects(the offerer
  and the acceptor) and passes session descriptions between them
  until the acceptor accepts the session*/
type AgentNegotiator func(Conn, Conn) error

//Conn represents a connection to an endpoint
type Conn interface {
	GetAddress() string
	IsClosed() bool
	io.ReadWriteCloser
}
