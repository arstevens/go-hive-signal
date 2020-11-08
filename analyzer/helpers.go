package analyzer

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

func calculateSplitPlacements(dspaceFrequencies map[string]int) []map[string]bool

func isValidMerge(fitScoreOne float64, fitScoreTwo float64) bool {
	distance := 1.0 - (fitScoreOne + fitScoreTwo)
	if distance < 0 {
		distance *= -1.0
	}
	return distance < MergeValidityLimit
}

func isValidSplit(fitScore float64, swarmSize int) bool {
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
