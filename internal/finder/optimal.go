package finder

var DefaultPreferredLoad = 1

/*OptimalSizeFinder implements analyzer.OptimalSizeFinder
and calculates the optimal size for a swarm based off of
its load parameter and the average preferred load of the
swarms members*/
type OptimalSizeFinder struct {
	tracker SwarmInfoTracker
}

//New creates a new instance of OptimalSizeFinder
func New(tracker SwarmInfoTracker) *OptimalSizeFinder {
	return &OptimalSizeFinder{
		tracker: tracker,
	}
}

//GetBestSize returns the optimal size for a given swarm
func (sf *OptimalSizeFinder) GetBestSize(swarmID string) int {
	swarmLoad := sf.tracker.GetLoad(swarmID)
	preferredLoadPer := sf.tracker.GetPreferredLoadPerMember(swarmID)
	if preferredLoadPer == 0 {
		preferredLoadPer = DefaultPreferredLoad
	}
	return swarmLoad / preferredLoadPer
}
