package analyzer

import (
	"sort"
	"sync"
	"time"
)

var DistancePollTime = time.Minute

func pollForNewDistances(tracker SwarmInfoTracker, optimalFinder OptimalSizeFinder,
	distances *swarmDistancesSlice, mutex *sync.Mutex) {
	for {
		time.Sleep(DistancePollTime)

		dataspaces := tracker.GetDataspaces()
		newDistances := make([]*swarmDistanceInfo, 0, len(dataspaces))
		for _, dataspace := range dataspaces {
			size := tracker.GetSize(dataspace)
			optimalSize := optimalFinder.GetBestSize(dataspace)

			distance := size - optimalSize
			newDistances = append(newDistances, &swarmDistanceInfo{
				dataspace: dataspace,
				distance:  distance,
			})
		}

		nd := swarmDistancesSlice(newDistances)
		sort.Sort(&nd)

		mutex.Lock()
		*distances = nd
		mutex.Unlock()
	}
}
