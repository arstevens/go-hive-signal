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

var FrequencyCalculationPeriod = time.Minute

var MergeValidityLimit = 0.3
var SplitSizeLimit = 50
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
		da.freqMutex.RLock()
		startDspaceLen, endDspaceLen := len(da.frequencies[startID]), len(da.frequencies[endID])
		da.freqMutex.RUnlock()
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
