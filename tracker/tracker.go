package tracker

import (
	"fmt"
	"math"
	"sync"
)

/*SwarmSizeTracker is an object that allows the tracking and
retrieving of swarm sizes*/
type SwarmSizeTracker struct {
	trackMap map[string]int
	mapMutex *sync.Mutex
}

//New creates a new instance of SwarmSizeTracker
func New() *SwarmSizeTracker {
	return &SwarmSizeTracker{
		trackMap: make(map[string]int),
		mapMutex: &sync.Mutex{},
	}
}

//GetSize returns the recorded size of the 'swarmID'
func (st *SwarmSizeTracker) GetSize(swarmID string) int {
	st.mapMutex.Lock()
	defer st.mapMutex.Unlock()
	if size, ok := st.trackMap[swarmID]; ok {
		return size
	}
	return 0
}

//Increment increments the size of 'swarmID' by one
func (st *SwarmSizeTracker) Increment(swarmID string) {
	st.mapMutex.Lock()
	if _, ok := st.trackMap[swarmID]; !ok {
		st.trackMap[swarmID] = 0
	}
	st.trackMap[swarmID]++
	st.mapMutex.Unlock()
}

//Decrement decrements the size of 'swarmID' by one
func (st *SwarmSizeTracker) Decrement(swarmID string) {
	st.mapMutex.Lock()
	size, ok := st.trackMap[swarmID]
	if !ok {
		st.trackMap[swarmID] = 0
	} else if size > 0 {
		st.trackMap[swarmID]--
	}
	st.mapMutex.Unlock()
}

//GetSmallest returns the smallest known swarm ID
func (st *SwarmSizeTracker) GetSmallest() (string, error) {
	minID := ""
	minSize := math.MaxInt32

	st.mapMutex.Lock()
	for swarmID, size := range st.trackMap {
		if size < minSize {
			minSize = size
			minID = swarmID
		}
	}
	st.mapMutex.Unlock()

	if minID == "" {
		return "", fmt.Errorf("No smallest swarm found in SwarmSizeTracker.GetSmallest()")
	}
	return minID, nil
}
