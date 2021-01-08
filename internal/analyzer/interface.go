package analyzer

/*SwarmInfoTracker describes an object that
has information on the number of members of
all known swarms as well as the load parameters
of each swarm*/
type SwarmInfoTracker interface {
	GetSize(string) int
	GetDataspaces() []string
}

/*OptimalSizeFinder provides an interface for
finding the optimal size of a swarm*/
type OptimalSizeFinder interface {
	GetBestSize(string) int
}
