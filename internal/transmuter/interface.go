package transmuter

import (
	"io"
)

/*SwarmSizeTracker describes an object that tracks
the number of members of each swarm*/
type SwarmSizeTracker interface {
	GetMostNeedy() (string, error)
}

/*SwarmMap describes an object that maps Swarm IDs
to the dataspaces they serve*/
type SwarmMap interface {
	GetSwarm(string) (SwarmManager, error)
}

/*SwarmAnalyzer describes an object that can make
recommendations for how to split/merge swarms*/
type SwarmAnalyzer interface {
	CalculateCandidates() ([]Candidate, error)
}

//Candidate describes a split or merge candidate
type Candidate interface {
	GetTransfererID() string
	GetTransfereeID() string
	GetTransferSize() int
}

type SwarmManager interface {
	AddEndpointConn(interface{}) error
	RemoveEndpointConn(interface{}) error
	AddEndpoint(string) error
	RemoveEndpoint(string) error
	GetEndpoints() []string
	SetID(string)
	io.Closer
}
