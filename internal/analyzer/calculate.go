package analyzer

import (
	"fmt"
	"sort"
	"sync"

	"github.com/arstevens/go-hive-signal/internal/transmuter"
)

/*OptimalLoadForSize must be set before DataRequestAnalyzer can
be used. It is a function that returns the optimal load to size
pairing for a swarm*/
var OptimalSizeForLoad func(size int) int = nil

type DataRequestAnalyzer struct {
	matchDistances swarmDistancesSlice
	dMutex         *sync.Mutex
	sizeTracker    SwarmInfoTracker
}

func New(sizeTracker SwarmInfoTracker) *DataRequestAnalyzer {
	analyzer := &DataRequestAnalyzer{
		matchDistances: swarmDistancesSlice(make([]*swarmDistanceInfo, 0)),
		dMutex:         &sync.Mutex{},
		sizeTracker:    sizeTracker,
	}
	go pollForNewDistances(sizeTracker, &analyzer.matchDistances, analyzer.dMutex)
	return analyzer
}

func (da *DataRequestAnalyzer) GetMostNeedy() (string, error) {
	da.dMutex.Lock()
	defer da.dMutex.Unlock()
	if da.matchDistances.Len() > 0 && da.matchDistances[0].distance < 0 {
		return da.matchDistances[0].dataspace, nil
	}
	return "", fmt.Errorf("Could not retrieve most needy swarm in DataRequestAnalyzer.GetMostNeedy()")
}

func (da *DataRequestAnalyzer) CalculateCandidates() ([]transmuter.Candidate, error) {
	da.dMutex.Lock()
	defer da.dMutex.Unlock()

	candidates := make([]transmuter.Candidate, 0)
	distances := da.matchDistances
	size := distances.Len()

	head, tail := distances[0], distances[size-1]
	for head.distance < 0 && tail.distance > 0 {
		tSize := calculateTransferSize(head, tail)
		candidates = append(candidates, &Candidate{
			transfererID: tail.dataspace,
			transfereeID: head.dataspace,
			transferSize: tSize,
		})

		head.distance += tSize
		tail.distance -= tSize
		adjustOrdering(&distances)
		head, tail = distances[0], distances[size-1]
	}
	return candidates, nil
}

func calculateTransferSize(neg *swarmDistanceInfo, pos *swarmDistanceInfo) int {
	negAbsDistance := -1 * neg.distance
	if negAbsDistance <= pos.distance {
		return negAbsDistance
	}
	return pos.distance
}

func adjustOrdering(container sort.Interface) {
	length := container.Len()
	for i := 0; i < length-2 && !container.Less(i, i+1); i++ {
		container.Swap(i, i+1)
	}

	for i := length - 1; i > 0 && !container.Less(i-1, i); i-- {
		container.Swap(i-1, i)
	}
}
