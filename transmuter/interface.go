package transmuter

import "github.com/arstevens/go-request/handle"

/*SwarmSizeTracker describes an object that tracks
the number of members of each swarm*/
type SwarmSizeTracker interface {
	GetSmallest() (string, error)
	Increment(string)
	Decrement(string)
}

/*SwarmMap describes an object that maps Swarm IDs
to the dataspaces they serve*/
type SwarmMap interface {
	RemoveSwarm(string) error
	/*Parameter: Swarm Dataspaces
	  Return Value: SwarmID, error*/
	AddSwarm([]string) (string, error)
	GetDataspaces(string) ([]string, error)
}

/*SwarmAnalyzer describes an object that can make
recommendations for how to split/merge swarms*/
type SwarmAnalyzer interface {
	CalculateCandidates() ([]Candidate, error)
}

//Candidate describes a split or merge candidate
type Candidate interface {
	IsSplit() bool
	GetSwarmIDs() []string
	/*Returns a set containing dataspaces that
	should be grouped together*/
	GetPlacementOne() map[string]bool
	GetPlacementTwo() map[string]bool
}

/*SwarmGateway describes an object that can perform
low-level swarm changing operations*/
type SwarmGateway interface {
	AddEndpoint(string, handle.Conn) error
	// Parameters: SwarmID, NewID_1, NewID_2
	Bisect(string, string, string) error
	// Parameters: SwarmID_1, SwarmID_2, NewID
	Stitch(string, string, string) error
}
