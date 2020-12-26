package tracker

import (
	"sync"
	"time"
)

var FrequencyCalculationPeriod = time.Minute

func calculateFrequencyOnInterval(trackerSize int, frequencies map[string]int, trackers map[string]*swarmLoadTracker,
	freqMutex *sync.Mutex, trackMutex *sync.Mutex) {
	for {
		time.Sleep(FrequencyCalculationPeriod)

		trackMutex.Lock()
		freqMutex.Lock()
		for dataspace, load := range frequencies {
			tracker, ok := trackers[dataspace]
			if !ok {
				tracker = newLoadTracker(trackerSize)
				trackers[dataspace] = tracker
			}
			tracker.AddFrequencyDatapoint(load)
			frequencies[dataspace] = 0
		}
		cleanup(frequencies, trackers)
		trackMutex.Unlock()
		freqMutex.Unlock()
	}
}

func cleanup(fmap map[string]int, tmap map[string]*swarmLoadTracker) {
	for dspace, tracker := range tmap {
		if tracker.CalculateAverageFrequency() == 0 {
			delete(tmap, dspace)
			delete(fmap, dspace)
		}
	}
}
