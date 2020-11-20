package manager

import "io"

/*SwarmGateway returns a connection object to an endpoint
in a designated swarm*/
type SwarmGateway interface {
	GetEndpoint(string) (interface{}, error)
}

/*AgentNegotiator takes in two connection objects(the offerer
  and the acceptor) and passes session descriptions between them
  until the acceptor accepts the session*/
type AgentNegotiator func(interface{}, interface{}) error

//Conn represents a connection to an endpoint
type Conn interface {
	io.ReadWriteCloser
}
