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

	totalSwarms := 5
	swarms := make([]string, totalSwarms)
	for i := 0; i < totalSwarms; i++ {
		swarms[i] = "/swarm/" + strconv.Itoa(i)
	}

	tracker := New()
	pollTime := time.Millisecond * 50
	totalReads := 10
	go func() {
		for i := 0; i < totalReads; i++ {
			time.Sleep(pollTime)
			smallestID, err := tracker.GetSmallest()
			fmt.Printf("Smallest Result: ID: %s, Size: %d, Error: %v\n", smallestID, tracker.GetSize(smallestID), err)
		}
	}()

	totalIncs := 300
	totalDecs := 60
	augPollTime := (pollTime * time.Duration(totalReads)) / time.Duration(totalIncs+totalDecs)

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < totalIncs; i++ {
		time.Sleep(augPollTime)
		swarmID := swarms[rand.Intn(len(swarms))]
		tracker.Increment(swarmID)
	}
	for i := 0; i < totalDecs; i++ {
		time.Sleep(augPollTime)
		swarmID := swarms[rand.Intn(len(swarms))]
		tracker.Decrement(swarmID)
	}
}
