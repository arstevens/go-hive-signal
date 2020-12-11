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

func (st *SwarmSizeTracker) SetSize(swarmID string, size int) {
	st.mapMutex.Lock()
	st.trackMap[swarmID] = size
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
