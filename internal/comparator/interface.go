package comparator

/*SwarmInfoTracker describes an object that has
information on the actual load parameter of a swarm
as well as the average preferred load for individual
members of the swarm*/
type SwarmInfoTracker interface {
	GetLoad(string) int
	GetDebriefData(string) interface{}
}
