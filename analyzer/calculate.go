package analyzer

import (
	"sort"
	"sync"
	"time"
)

/*OptimalLoadForSize must be set before DataRequestAnalyzer can
be used. It is a function that returns the optimal load to size
pairing for a swarm*/
var OptimalLoadForSize func(size int) int = nil

/*IncrementModifier must be set before DataRequestAnalyzer can be
used. It is a modifier that changes the way the load/size fit is
calculated*/
var IncrementModifier int = 0

/*FrequencyCalculationPeriod sets how often a new frequency data
point will be calculated and added to the record*/
var FrequencyCalculationPeriod = time.Minute

/*MergeValidityLimit is the maximum distance from a size/load fit of
1.0 that a merge candidate can be to be considered valid. Since a score
closed to 0.0 means a swarm is small and has too much load and 1.0 means
a swarm is big and has too little load it follows that an optimal sum of
fit scores should be one since a small swarm with lots of data should be
paired with a large swarm with little data*/
var MergeValidityLimit = 0.3

/*SplitSizeLimit is the minimum size of a swarm in order for it to
be considered eligable for a split*/
var SplitSizeLimit = 50

/*SplitValidityLimit is the minimum distance from a perfect fit(0.5)
that a swarm must be in order to be a valid candidate for a split*/
var SplitValidityLimit = 0.3

type DataRequestAnalyzer struct {
	/*maps Swarm IDs to a map of Dataspace IDs that
	  contain frequency information */
	frequencies map[string]map[string]int
	freqMutex   *sync.RWMutex
	trackers    map[string]*swarmTracker
	trackMutex  *sync.RWMutex
	sizeTracker SwarmSizeTracker
}

func New(sizeTracker SwarmSizeTracker) *DataRequestAnalyzer {
	frequencies := make(map[string]map[string]int)
	trackers := make(map[string]*swarmTracker)
	fMutex := &sync.RWMutex{}
	tMutex := &sync.RWMutex{}

	go calculateFrequencyOnInterval(frequencies, trackers, fMutex, tMutex)
	return &DataRequestAnalyzer{
		frequencies: frequencies,
		trackers:    trackers,
		freqMutex:   fMutex,
		trackMutex:  tMutex,
		sizeTracker: sizeTracker,
	}
}

func (da *DataRequestAnalyzer) CalculateCandidates() ([]Candidate, error) {
	candidates := make([]Candidate, 0)
	da.trackMutex.RLock()
	fits := calculateSwarmFits(da.trackers, da.sizeTracker)
	fitList := swarmInfoList{Infos: fits}
	sort.Sort(&fitList)

	start := 0
	end := len(fits) - 1
	for start <= end {
		startID, endID := fits[start].SwarmID, fits[end].SwarmID
		startFit, endFit := fits[start].FitScore, fits[end].FitScore
		startSize, endSize := fits[start].SwarmSize, fits[end].SwarmSize
		startDspaceLen := da.trackers[startID].TotalActiveDataspaces()
		endDspaceLen := da.trackers[endID].TotalActiveDataspaces()
		if isValidSplit(startFit, startSize, startDspaceLen) {
			tracker := da.trackers[fits[start].SwarmID]
			candidate := createSplitCandidate(&fits[start], tracker)
			candidates = append(candidates, candidate)
		} else if isValidSplit(endFit, endSize, endDspaceLen) {
			tracker := da.trackers[fits[end].SwarmID]
			candidate := createSplitCandidate(&fits[end], tracker)
			candidates = append(candidates, candidate)
		} else if start != end && isValidMerge(startFit, endFit) {
			candidate := Candidate{
				isSplit: false,
				swarms:  []string{fits[start].SwarmID, fits[end].SwarmID},
			}
			candidates = append(candidates, candidate)
		}
		start++
		end--
	}
	da.trackMutex.RUnlock()
	return candidates, nil
}

func (da *DataRequestAnalyzer) IncrementFrequencyCounter(swarmID string, dspaceID string) {
	da.freqMutex.Lock()
	dataspaces, ok := da.frequencies[swarmID]
	if !ok {
		dataspaces = make(map[string]int)
		da.frequencies[swarmID] = dataspaces
	}
	da.freqMutex.Unlock()

	if _, ok := dataspaces[dspaceID]; !ok {
		dataspaces[dspaceID] = 0
	}
	dataspaces[dspaceID]++
}

func calculateFrequencyOnInterval(frequencies map[string]map[string]int, trackers map[string]*swarmTracker,
	freqMutex *sync.RWMutex, trackMutex *sync.RWMutex) {
	for {
		time.Sleep(FrequencyCalculationPeriod)

		trackMutex.Lock()
		freqMutex.Lock()
		for swarmID, dataspaces := range frequencies {
			tracker, ok := trackers[swarmID]
			if !ok {
				tracker = newTracker()
				trackers[swarmID] = tracker
			}
			for dspaceID, dFreq := range dataspaces {
				tracker.AddFrequencyDatapoint(dspaceID, dFreq)
				dataspaces[dspaceID] = 0
			}
		}
		cleanup(trackers, frequencies)
		trackMutex.Unlock()
		freqMutex.Unlock()
	}
}

func cleanup(tmap map[string]*swarmTracker, fmap map[string]map[string]int) {
	for key, swarm := range tmap {
		if swarm.CalculateFrequency() == 0 {
			delete(tmap, key)
			delete(fmap, key)
		} else {
			swarm.Cleanup()
		}
	}
}
