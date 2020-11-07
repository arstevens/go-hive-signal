package transmuter

import "github.com/arstevens/go-request/handle"

type SwarmSizeTracker interface {
	GetSmallest() (string, error)
	Increment(string)
	Decrement(string)
}

type SwarmMap interface {
	RemoveSwarm(string) error
	/*Parameter: Swarm Dataspaces
	  Return Value: SwarmID, error*/
	AddSwarm([]string) (string, error)
	GetDataspaces(string) ([]string, error)
}

type SwarmAnalyzer interface {
	GetCandidates() ([]Candidate, error)
}

type Candidate interface {
	IsSplit() bool
	GetSwarmIDs() []string
	GetPlacementOne() map[string]bool
	GetPlacementTwo() map[string]bool
}

type SwarmGateway interface {
	AddEndpoint(string, handle.Conn) error
	// Parameters: SwarmID, NewID_1, NewID_2
	Bisect(string, string, string) error
	// Parameters: SwarmID_1, SwarmID_2, NewID
	Stitch(string, string, string) error
}
