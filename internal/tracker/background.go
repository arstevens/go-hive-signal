package tracker

import (
	"sync"
	"time"
)

var FrequencyCalculationPeriod = time.Minute

func calculateFrequencyOnInterval(trackerSize int, frequencies map[string]*loadEntry, trackers map[string]*swarmLoadTracker,
	freqMutex *sync.RWMutex, trackMutex *sync.Mutex) {
	for {
		time.Sleep(FrequencyCalculationPeriod)

		trackMutex.Lock()
		freqMutex.Lock()
		for dataspace, entry := range frequencies {
			tracker, ok := trackers[dataspace]
			if !ok {
				tracker = newLoadTracker(trackerSize)
				trackers[dataspace] = tracker
			}
			entry.mutex.Lock()
			tracker.AddFrequencyDatapoint(entry.load)
			entry.load = 0
			entry.mutex.Unlock()
		}
		cleanup(frequencies, trackers)
		trackMutex.Unlock()
		freqMutex.Unlock()
	}
}

func cleanup(fmap map[string]*loadEntry, tmap map[string]*swarmLoadTracker) {
	for dspace, tracker := range tmap {
		if tracker.CalculateAverageFrequency() == 0 {
			delete(tmap, dspace)
			delete(fmap, dspace)
		}
	}
}
