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

func New() *DataRequestAnalyzer {
	frequencies := make(map[string]map[string]int)
	trackers := make(map[string]*swarmTracker)
	fMutex := &sync.RWMutex{}
	tMutex := &sync.RWMutex{}

	calculateFrequencyOnInterval(frequencies, trackers, fMutex, tMutex)
	return &DataRequestAnalyzer{
		frequencies: frequencies,
		trackers:    trackers,
		freqMutex:   fMutex,
		trackMutex:  tMutex,
	}
}

func (da *DataRequestAnalyzer) CalculateCandidates() ([]Candidate, error) {
	candidates := make([]Candidate, 0)
	fits := calculateSwarmFits(da.trackers, da.sizeTracker)
	fitList := swarmInfoList{Infos: fits}
	sort.Sort(&fitList)

	start := 0
	end := len(fits) - 1
	for start <= end {
		if isValidSplit(fits[start].FitScore, fits[start].SwarmSize) {
			tracker := da.trackers[fits[start].SwarmID]
			candidate := createSplitCandidate(&fits[start], tracker)
			candidates = append(candidates, candidate)
		} else if isValidSplit(fits[end].FitScore, fits[end].SwarmSize) {
			tracker := da.trackers[fits[end].SwarmID]
			candidate := createSplitCandidate(&fits[end], tracker)
			candidates = append(candidates, candidate)
		} else if start != end && isValidMerge(fits[start].FitScore, fits[end].FitScore) {
			candidate := Candidate{
				isSplit: false,
				swarms:  []string{fits[start].SwarmID, fits[end].SwarmID},
			}
			candidates = append(candidates, candidate)
		}
		start++
		end--
	}
	return candidates, nil
}

func (da *DataRequestAnalyzer) IncrementFrequencyCounter(swarmID string, dspaceID string) {
	dataspaces, ok := da.frequencies[swarmID]
	if !ok {
		dataspaces = make(map[string]int)
		da.frequencies[swarmID] = dataspaces
	}

	if _, ok := dataspaces[dspaceID]; !ok {
		dataspaces[dspaceID] = 0
	}
	dataspaces[dspaceID]++
}

func calculateFrequencyOnInterval(frequencies map[string]map[string]int, trackers map[string]*swarmTracker,
	freqMutex *sync.RWMutex, trackMutex *sync.RWMutex) {
	for {
		time.Sleep(FrequencyCalculationPeriod)

		for swarmID, dataspaces := range frequencies {
			frequency := 0
			tracker, ok := trackers[swarmID]
			if !ok {
				tracker = newTracker()
				trackers[swarmID] = tracker
			}
			for dspaceID, dFreq := range dataspaces {
				tracker.AddDataspaceFrequencyDatapoint(dspaceID, dFreq)
				frequency += dFreq
				dataspaces[dspaceID] = 0
			}
			tracker.AddFrequencyDatapoint(frequency)
		}
	}
}
