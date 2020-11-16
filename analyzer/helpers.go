package analyzer

import (
	"sort"
)

func createSplitCandidate(fit *swarmInfo, tracker *swarmTracker) Candidate {
	placements := calculateSplitPlacements(tracker.CalculateDataspaceFrequencies())
	candidate := Candidate{
		isSplit:      true,
		swarms:       []string{fit.SwarmID},
		placementOne: placements[0],
		placementTwo: placements[1],
	}
	return candidate
}

func calculateSplitPlacements(dspaceFrequencies map[string]int) []map[string]bool {
	sortedDataspaces := sortDataspacesByFrequency(dspaceFrequencies)

	setOne := make(map[string]bool)
	setTwo := make(map[string]bool)

	idx := len(sortedDataspaces) - 1
	cumulativeLoadOne := dspaceFrequencies[sortedDataspaces[idx]]
	setOne[sortedDataspaces[idx]] = true
	idx--
	cumulativeLoadTwo := dspaceFrequencies[sortedDataspaces[idx]]
	setTwo[sortedDataspaces[idx]] = true
	idx--

	for idx >= 0 {
		id := sortedDataspaces[idx]
		load := dspaceFrequencies[id]
		if cumulativeLoadOne < cumulativeLoadTwo {
			setOne[id] = true
			cumulativeLoadOne += load
		} else {
			setTwo[id] = true
			cumulativeLoadTwo += load
		}
		idx--
	}
	return []map[string]bool{setOne, setTwo}
}

func sortDataspacesByFrequency(dspaceFrequencies map[string]int) []string {
	dataspaces := make([]string, len(dspaceFrequencies))
	i := 0
	for key, _ := range dspaceFrequencies {
		dataspaces[i] = key
		i++
	}
	sort.Sort(&dataspaceFrequencyPair{dataspaces, dspaceFrequencies})
	return dataspaces
}

type dataspaceFrequencyPair struct {
	dataspaces  []string
	frequencies map[string]int
}

func (dp *dataspaceFrequencyPair) Len() int { return len(dp.dataspaces) }
func (dp *dataspaceFrequencyPair) Less(i, j int) bool {
	return dp.frequencies[dp.dataspaces[i]] < dp.frequencies[dp.dataspaces[i]]
}
func (dp *dataspaceFrequencyPair) Swap(i, j int) {
	dp.dataspaces[i], dp.dataspaces[j] = dp.dataspaces[j], dp.dataspaces[i]
}

func isValidMerge(fitScoreOne float64, fitScoreTwo float64) bool {
	distance := 1.0 - (fitScoreOne + fitScoreTwo)
	if distance < 0 {
		distance *= -1.0
	}
	return distance < MergeValidityLimit
}

func isValidSplit(fitScore float64, swarmSize int, totalDspaces int) bool {
	if totalDspaces < 2 {
		return false
	}

	sizeValidity := swarmSize >= SplitSizeLimit
	if fitScore < 0.5 {
		fitScore = 1.0 - fitScore
	}
	fitValidity := (fitScore - 0.5) > SplitValidityLimit
	return sizeValidity && fitValidity
}

func calculateSwarmFits(trackers map[string]*swarmTracker, sizeTracker SwarmSizeTracker) []swarmInfo {
	fits := make([]swarmInfo, 0)
	for swarmID, tracker := range trackers {
		swarmSize := sizeTracker.GetSize(swarmID)
		loadParameter := tracker.CalculateFrequency()

		fitScore := compareSizeToLoad(swarmSize, loadParameter)
		sInfo := swarmInfo{SwarmID: swarmID, SwarmSize: swarmSize, FitScore: fitScore}
		fits = append(fits, sInfo)
	}
	return fits
}

func compareSizeToLoad(swarmSize int, loadParameter int) float64 {
	distance := OptimalLoadForSize(swarmSize) - loadParameter
	if distance == 0 {
		return 0.5
	}

	neg := false
	if distance < 0 {
		distance *= -1
		neg = true
	}

	metric := float64(distance) / float64(distance+IncrementModifier)
	if neg {
		return 1.0 - metric
	}
	return metric
}
