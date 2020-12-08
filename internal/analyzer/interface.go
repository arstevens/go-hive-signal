package analyzer

/*SwarmSizeTracker describes an object that
has information on the number of members of
all known swarms*/
type SwarmSizeTracker interface {
	GetSize(string) int
}
