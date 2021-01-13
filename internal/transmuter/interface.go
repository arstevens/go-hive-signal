package transmuter

import (
	"io"
)

/*SwarmMap describes an object that maps Swarm IDs
to the dataspaces they serve*/
type SwarmMap interface {
	GetSwarm(string) (interface{}, error)
}

/*SwarmAnalyzer describes an object that can make
recommendations for how to split/merge swarms*/
type SwarmAnalyzer interface {
	CalculateCandidates() ([]Candidate, error)
	GetMostNeedy() (string, error)
}

//Candidate describes a split or merge candidate
type Candidate interface {
	GetTransfererID() string
	GetTransfereeID() string
	GetTransferSize() int
}

type SwarmManager interface {
	AddEndpoint(interface{}) error
	Transfer(int, SwarmManager) error
	io.Closer
}
