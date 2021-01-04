package transmuter

import (
	"fmt"
	"log"
	"time"
)

const (
	splitSwarmFailFormat     = "Failed to split swarm map in SwarmTransmuter: %v"
	mergeSwarmFailFormat     = "Failed to merge swarm map in SwarmTransmuter: %v"
	transmuteSwarmFailFormat = "Failed to transmute swarms in SwarmTransmuter: %v"
)

const (
	bisectKey = iota
	mergeKey
)

var PollPeriod = time.Minute

func pollForTransmutation(swarmMap SwarmMap, analyzer SwarmAnalyzer) {
	for {
		time.Sleep(PollPeriod)
		candidates, err := analyzer.CalculateCandidates()
		if err != nil {
			log.Println(err)
			return
		}

		if len(candidates) > 0 {
			err := transmuteSwarms(swarmMap, candidates)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func transmuteSwarms(swarmMap SwarmMap, candidates []Candidate) error {
	for _, candidate := range candidates {
		transfererID := candidate.GetTransfererID()
		transfereeID := candidate.GetTransfereeID()
		transferSize := candidate.GetTransferSize()

		t, err := swarmMap.GetSwarm(transfererID)
		if err != nil {
			return fmt.Errorf("Failed to retrieve transferer SwarmManager of dataspace %s in SwarmTransmuter daemon: %v", transfererID, err)
		}
		transferer := t.(SwarmManager)

		t, err = swarmMap.GetSwarm(transfereeID)
		if err != nil {
			return fmt.Errorf("Failed to retrieve transferee SwarmManager of dataspace %s in SwarmTransmuter daemon: %v", transfererID, err)
		}
		transferee := t.(SwarmManager)

		endpoints := transferer.GetEndpoints()
		for i := 0; i < len(endpoints) && i < transferSize; i++ {
			err = transferee.TakeEndpoint(endpoints[i])
			if err != nil {
				return fmt.Errorf("Failed to transfer endpoint %s in SwarmTransmuter daemon: %v", endpoints[i], err)
			}
			transferer.DropEndpoint(endpoints[i])
		}
	}
	return nil
}
