package tracker

import (
	"sync"
)

/*SwarmSizeTracker is an object that allows the tracking and
retrieving of swarm sizes*/
type SwarmSizeTracker struct {
	loadMap       map[string]int
	lMapMutex     *sync.Mutex
	trackers      map[string]*swarmLoadTracker
	trackersMutex *sync.Mutex
	trackMap      map[string]int
	tMapMutex     *sync.Mutex
}

//New creates a new instance of SwarmSizeTracker
func New(historyLength int) *SwarmSizeTracker {
	tracker := &SwarmSizeTracker{
		loadMap:       make(map[string]int),
		lMapMutex:     &sync.Mutex{},
		trackers:      make(map[string]*swarmLoadTracker),
		trackersMutex: &sync.Mutex{},
		trackMap:      make(map[string]int),
		tMapMutex:     &sync.Mutex{},
	}
	go calculateFrequencyOnInterval(historyLength, tracker.loadMap, tracker.trackers,
		tracker.lMapMutex, tracker.trackersMutex)
	return tracker
}

func (st *SwarmSizeTracker) IncrementFrequencyCounter(dataspace string) {
	st.lMapMutex.Lock()
	if _, ok := st.loadMap[dataspace]; !ok {
		st.loadMap[dataspace] = 0
	}
	st.loadMap[dataspace]++
	st.lMapMutex.Unlock()
}

func (st *SwarmSizeTracker) GetLoad(dataspace string) int {
	st.trackersMutex.Lock()
	defer st.trackersMutex.Unlock()
	if tracker, ok := st.trackers[dataspace]; ok {
		return tracker.CalculateAverageFrequency()
	}
	return 0
}

//GetSize returns the recorded size of the 'swarmID'
func (st *SwarmSizeTracker) GetSize(swarmID string) int {
	st.tMapMutex.Lock()
	defer st.tMapMutex.Unlock()
	if size, ok := st.trackMap[swarmID]; ok {
		return size
	}
	return 0
}

func (st *SwarmSizeTracker) SetSize(swarmID string, size int) {
	st.tMapMutex.Lock()
	st.trackMap[swarmID] = size
	st.tMapMutex.Unlock()
}
