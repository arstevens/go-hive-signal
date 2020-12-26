package tracker

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func TestTracker(t *testing.T) {
	fmt.Printf("---------------TRACKER TEST------------------\n")

	totalSwarms := 10
	swarms := make([]string, totalSwarms)
	for i := 0; i < totalSwarms; i++ {
		swarms[i] = "/dataspace/" + strconv.Itoa(i)
	}

	FrequencyCalculationPeriod = time.Millisecond * 10
	tracker := New(0)
	totalOps := 300
	freqCountPer := 1000
	pollTime := time.Millisecond * 3

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < totalOps; i++ {
		time.Sleep(pollTime)

		dspace := swarms[rand.Intn(totalSwarms)]
		tracker.SetSize(dspace, tracker.GetSize(dspace)+1)

		count := rand.Intn(freqCountPer)
		for j := 0; j < count; j++ {
			tracker.IncrementFrequencyCounter(dspace)
		}
	}

	for _, dspace := range swarms {
		fmt.Printf("Swarm %s : Size %d : Load %d\n", dspace, tracker.GetSize(dspace), tracker.GetLoad(dspace))
	}
}
