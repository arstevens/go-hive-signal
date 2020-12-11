package transmuter

import "github.com/arstevens/go-request/handle"

/*SwarmSizeTracker describes an object that tracks
the number of members of each swarm*/
type SwarmSizeTracker interface {
	GetSmallest() (string, error)
}

/*SwarmMap describes an object that maps Swarm IDs
to the dataspaces they serve*/
type SwarmMap interface {
	RemoveSwarm(string) error
	/*Parameter: Swarm Dataspaces
	  Return Value: SwarmID, error*/
	AddSwarm(SwarmManager, []string) (string, error)
	GetDataspaces(string) ([]string, error)
	GetSwarmByID(string) (SwarmManager, error)
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

type SwarmManager interface {
	AddEndpoint(handle.Conn) error
	RemoveEndpoint(handle.Conn) error
	// Parameters New ID 1, New ID 2
	Bisect() (SwarmManager, error)
	// Parameters Swarm to merge, New ID
	Stitch(SwarmManager) error
	Destroy() error
}
