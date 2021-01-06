package tracker

import (
	"sync"
)

type loadEntry struct {
	load  int
	mutex *sync.Mutex
}

/*SwarmInfoTracker is an object that allows the tracking and
retrieving of swarm sizes*/
type SwarmInfoTracker struct {
	loadMap       map[string]*loadEntry
	loadMutex     *sync.RWMutex
	trackers      map[string]*swarmLoadTracker
	trackersMutex *sync.Mutex
	sizeMap       map[string]int
	sizeMutex     *sync.RWMutex
}

//New creates a new instance of SwarmInfoTracker
func New(historyLength int) *SwarmInfoTracker {
	tracker := &SwarmInfoTracker{
		loadMap:       make(map[string]*loadEntry),
		loadMutex:     &sync.RWMutex{},
		trackers:      make(map[string]*swarmLoadTracker),
		trackersMutex: &sync.Mutex{},
		sizeMap:       make(map[string]int),
		sizeMutex:     &sync.RWMutex{},
	}
	go calculateFrequencyOnInterval(historyLength, tracker.loadMap, tracker.trackers,
		tracker.loadMutex, tracker.trackersMutex)
	return tracker
}

func (st *SwarmInfoTracker) IncrementFrequencyCounter(dataspace string) {
	st.loadMutex.RLock()
	var ok bool
	var entry *loadEntry
	if entry, ok = st.loadMap[dataspace]; !ok {
		st.loadMutex.RUnlock()
		st.loadMutex.Lock()
		st.loadMap[dataspace] = &loadEntry{load: 0, mutex: &sync.Mutex{}}
		entry = st.loadMap[dataspace]
		st.loadMutex.Unlock()
		st.loadMutex.RLock()
	}

	entry.mutex.Lock()
	entry.load++
	entry.mutex.Unlock()
	st.loadMutex.RUnlock()
}

func (st *SwarmInfoTracker) GetLoad(dataspace string) int {
	st.trackersMutex.Lock()
	defer st.trackersMutex.Unlock()
	if tracker, ok := st.trackers[dataspace]; ok {
		return tracker.CalculateAverageFrequency()
	}
	return 0
}

func (st *SwarmInfoTracker) GetDataspaces() []string {
	st.sizeMutex.RLock()
	dspaces := make([]string, 0, len(st.sizeMap))
	for dspace, _ := range st.sizeMap {
		dspaces = append(dspaces, dspace)
	}
	st.sizeMutex.RUnlock()
	return dspaces
}

//GetSize returns the recorded size of the 'swarmID'
func (st *SwarmInfoTracker) GetSize(swarmID string) int {
	st.sizeMutex.RLock()
	defer st.sizeMutex.RUnlock()
	if size, ok := st.sizeMap[swarmID]; ok {
		return size
	}
	return 0
}

func (st *SwarmInfoTracker) SetSize(swarmID string, size int) {
	st.sizeMutex.Lock()
	st.sizeMap[swarmID] = size
	st.sizeMutex.Unlock()
}

func (st *SwarmInfoTracker) Delete(swarmID string) {
	st.sizeMutex.Lock()
	delete(st.sizeMap, swarmID)
	st.sizeMutex.Unlock()
}
