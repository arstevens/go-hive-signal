package analyzer

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func TestIntegration(t *testing.T) {
	fmt.Printf("---------------INTEGRATION TEST------------------\n")

	totalSwarms := 20
	totalDspacesPer := 10
	dpointRange := 100
	totalDpointsPer := 50
	sizeRange := 100

	rand.Seed(time.Now().UnixNano())
	swarmIDs := make([]string, 0)
	sDspaceIDs := make(map[string][]string)
	swarms := make(map[string]*swarmTracker)
	for i := 0; i < totalSwarms; i++ {
		sid := "/swarm/" + strconv.Itoa(i)
		swarmIDs = append(swarmIDs, sid)
		swarms[sid] = newTracker()
		totalDspaces := rand.Intn(totalDspacesPer)
		sDspaceIDs[sid] = []string{}
		for j := 0; j < totalDspaces; j++ {
			did := "/dataspace/" + strconv.Itoa(i*totalDspacesPer+j)
			sDspaceIDs[sid] = append(sDspaceIDs[sid], did)
			totalDpoints := rand.Intn(totalDpointsPer)
			for k := 0; k < totalDpoints; k++ {
				dpoint := rand.Intn(dpointRange)
				swarms[sid].AddFrequencyDatapoint(did, dpoint)
			}
		}
	}

	sizes := make(map[string]int)
	for id, _ := range swarms {
		sizes[id] = rand.Intn(sizeRange)
	}
	sizeTracker := TestSwarmSizeTracker{sizes: sizes}

	OptimalLoadForSize = func(size int) int {
		return size
	}
	IncrementModifier = 20
	SplitValidityLimit = 0.1
	FrequencyCalculationPeriod = time.Millisecond * 10
	analyzer := New(&sizeTracker)
	analyzer.trackers = swarms

	totalIncs := 1000
	incFreq := time.Millisecond
	go func() {
		for i := 0; i < totalIncs; i++ {
			time.Sleep(incFreq)
			randSwarm := swarmIDs[rand.Intn(len(swarmIDs))]
			totalDspaces := len(sDspaceIDs[randSwarm])
			if totalDspaces > 0 {
				randDspace := sDspaceIDs[randSwarm][rand.Intn(totalDspaces)]
				analyzer.IncrementFrequencyCounter(randSwarm, randDspace)
			}
		}
	}()

	pollTime := time.Millisecond
	totalPolls := incFreq * time.Duration(totalIncs) / pollTime
	for i := 0; i < int(totalPolls); i++ {
		time.Sleep(pollTime)
		_, err := analyzer.CalculateCandidates()
		if err != nil {
			panic(err)
		}
	}
	fmt.Println("No deadlocks or race conditions detected")
}

func TestSwarmFit(t *testing.T) {
	fmt.Printf("---------------FITNESS CALCULATION TEST------------------\n")
	totalSwarms := 20
	totalDspacesPer := 10
	dpointRange := 100
	totalDpointsPer := 50
	sizeRange := 100

	rand.Seed(time.Now().UnixNano())
	swarms := make(map[string]*swarmTracker)
	for i := 0; i < totalSwarms; i++ {
		sid := "/swarm/" + strconv.Itoa(i)
		swarms[sid] = newTracker()
		totalDspaces := rand.Intn(totalDspacesPer)
		for j := 0; j < totalDspaces; j++ {
			did := "/dataspace/" + strconv.Itoa(i*totalDspacesPer+j)
			totalDpoints := rand.Intn(totalDpointsPer)
			for k := 0; k < totalDpoints; k++ {
				dpoint := rand.Intn(dpointRange)
				swarms[sid].AddFrequencyDatapoint(did, dpoint)
			}
		}
	}

	sizes := make(map[string]int)
	for id, _ := range swarms {
		sizes[id] = rand.Intn(sizeRange)
	}
	sizeTracker := TestSwarmSizeTracker{sizes: sizes}

	OptimalLoadForSize = func(size int) int {
		return size
	}
	IncrementModifier = 20

	fits := calculateSwarmFits(swarms, &sizeTracker)
	for _, fit := range fits {
		printSwarmInfo(fit)
		totalDspaces := len(swarms[fit.SwarmID].dspaceFrequencyHistory)
		fmt.Printf("IsValidSplit: %t\n", isValidSplit(fit.FitScore, fit.SwarmSize, totalDspaces))
	}
}

func printSwarmInfo(sInfo swarmInfo) {
	fmt.Println("{")
	fmt.Printf("\tID: %s\n", sInfo.SwarmID)
	fmt.Printf("\tSize: %d\n", sInfo.SwarmSize)
	fmt.Printf("\tFit: %f\n", sInfo.FitScore)
	fmt.Println("}")
}

func TestCandidateCalculation(t *testing.T) {
	fmt.Printf("---------------CANDIDATE CALCULATION TEST------------------\n")
	totalSwarms := 20
	totalDspacesPer := 10
	dpointRange := 100
	totalDpointsPer := 50
	sizeRange := 100

	rand.Seed(time.Now().UnixNano())
	swarms := make(map[string]*swarmTracker)
	for i := 0; i < totalSwarms; i++ {
		sid := "/swarm/" + strconv.Itoa(i)
		swarms[sid] = newTracker()
		totalDspaces := rand.Intn(totalDspacesPer)
		for j := 0; j < totalDspaces; j++ {
			did := "/dataspace/" + strconv.Itoa(i*totalDspacesPer+j)
			totalDpoints := rand.Intn(totalDpointsPer)
			for k := 0; k < totalDpoints; k++ {
				dpoint := rand.Intn(dpointRange)
				swarms[sid].AddFrequencyDatapoint(did, dpoint)
			}
		}
	}

	sizes := make(map[string]int)
	for id, _ := range swarms {
		sizes[id] = rand.Intn(sizeRange)
	}
	sizeTracker := TestSwarmSizeTracker{sizes: sizes}

	OptimalLoadForSize = func(size int) int {
		return size
	}
	IncrementModifier = 20
	SplitValidityLimit = 0.1
	analyzer := New(&sizeTracker)
	analyzer.trackers = swarms
	candidates, err := analyzer.CalculateCandidates()
	if err != nil {
		panic(err)
	}

	for _, candidate := range candidates {
		printCandidate(candidate)
	}
}

func printCandidate(c Candidate) {
	fmt.Println("{")
	fmt.Printf("\tIsSplit -> (%t)\n", c.isSplit)
	fmt.Printf("\tSwarms -> ")
	for idx, swarm := range c.swarms {
		prefix := ", "
		if idx == 0 {
			prefix = ""
		}
		fmt.Printf("%s%s", prefix, swarm)
	}
	fmt.Println()
	fmt.Printf("\tPlacementOne -> {\n")
	for id, _ := range c.placementOne {
		fmt.Printf("\t\t%s\n", id)
	}
	fmt.Printf("\t}\n")
	fmt.Printf("\tPlacementTwo -> {\n")
	for id, _ := range c.placementTwo {
		fmt.Printf("\t\t%s\n", id)
	}
	fmt.Printf("\t}\n")
	fmt.Printf("}\n")
}

func TestFrequencyCalculation(t *testing.T) {
	fmt.Printf("---------------FREQUENCY CALCULATION TEST------------------\n")
	FrequencyCalculationPeriod = time.Millisecond * 10
	sizeTracker := TestSwarmSizeTracker{sizes: make(map[string]int)}
	dra := New(&sizeTracker)

	totalMilis := 5
	incFreq := time.Millisecond * time.Duration(totalMilis)
	totalIncs := 100

	totalSwarms := 5
	swarms := make([]string, totalSwarms)
	for i := 0; i < totalSwarms; i++ {
		swarms[i] = "/swarm/" + strconv.Itoa(i)
	}

	for i := 0; i < totalIncs; i++ {
		swarmPick := rand.Intn(totalSwarms)
		dspace := "/dataspace/" + strconv.Itoa(swarmPick)
		for j := 0; j < 100; j++ {
			dra.IncrementFrequencyCounter(swarms[swarmPick], dspace)
		}
		time.Sleep(incFreq)
		if i == totalIncs/4 || i == totalIncs/2 || i == 3*totalIncs/4 || i == totalIncs-1 {
			fmt.Printf("----------------\nStatus at %d ms\n----------------\n", i*totalMilis)
			printTrackers(dra)
		}
	}

	time.Sleep(FrequencyCalculationPeriod * time.Duration(FrequencyAveragingWidth+5))
	fmt.Printf("----------------\nStatus post kill\n-----------------\n")
	printTrackers(dra)
}

func printTrackers(dra *DataRequestAnalyzer) {
	dra.trackMutex.Lock()
	entered := false
	for id, tracker := range dra.trackers {
		entered = true
		fmt.Printf("Swarm %s has frequency %d\n", id, tracker.CalculateFrequency())
	}
	if !entered {
		fmt.Printf("NO TRACKING DATA\n")
	}
	dra.trackMutex.Unlock()
}

func TestSwarmTracker(t *testing.T) {
	fmt.Printf("---------------SWARM TRACKER TEST------------------\n")
	FrequencyAveragingWidth = 5
	tracker := newTracker()

	totalDspaces := 10
	totalDspaceDatapoints := 10

	dataspaces := make([]string, totalDspaces)
	for i := 0; i < totalDspaces; i++ {
		dataspaces[i] = "/dataspace/" + strconv.Itoa(i)
	}

	fmt.Printf("frequency calculations\n----------------------\n")
	for i := 0; i < totalDspaces; i++ {
		dspace := dataspaces[i]
		for j := 0; j < totalDspaceDatapoints; j++ {
			tracker.AddFrequencyDatapoint(dspace, rand.Intn(100))
			averages := tracker.CalculateDataspaceFrequencies()
			fmt.Printf("Average Frequency for (%s): %d\n", dspace, averages[dspace])
		}
		fmt.Printf("Total Average Frequency: %d\n", tracker.CalculateFrequency())
	}

	fmt.Printf("------------------\ngarbage collection\n------------------\n")
	fmt.Printf("Pre-removal frequency: %d\n", tracker.CalculateFrequency())
	for i := 0; i < totalDspaces/2; i++ {
		dspace := dataspaces[i]
		for j := 0; j < FrequencyAveragingWidth; j++ {
			tracker.AddFrequencyDatapoint(dspace, 0)
		}
	}
	tracker.Cleanup()
	fmt.Printf("Post-removal frequency: %d\n", tracker.CalculateFrequency())
}

type TestSwarmSizeTracker struct {
	sizes map[string]int
}

func (tt *TestSwarmSizeTracker) GetSize(id string) int {
	size, ok := tt.sizes[id]
	if !ok {
		return 0
	}
	return size
}
