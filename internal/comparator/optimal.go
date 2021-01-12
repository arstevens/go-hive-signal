package comparator

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
	debriefData := sf.tracker.GetDebriefData(swarmID)

	preferredLoadPer := DefaultPreferredLoad
	if debriefData != nil && debriefData.(int) != 0 {
		preferredLoadPer = debriefData.(int)
	}
	return swarmLoad / preferredLoadPer
}
