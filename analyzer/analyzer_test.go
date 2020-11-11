package analyzer

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

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
