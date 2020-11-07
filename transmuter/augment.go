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

var PollPeriod = time.Second

func pollForTransmutation(swarmMap SwarmMap, gateway SwarmGateway, analyzer SwarmAnalyzer) {
	for {
		time.Sleep(PollPeriod)
		candidates, err := analyzer.GetCandidates()
		if err != nil {
			return
		}

		if len(candidates) > 0 {
			transmuteInstructions, err := transmuteSwarmMap(swarmMap, candidates)
			if err != nil {
				log.Println(err)
			}
			err = transmuteSwarms(gateway, transmuteInstructions)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func transmuteSwarms(gateway SwarmGateway, swarms map[int][][]string) error {
	bisects := swarms[bisectKey]
	for _, candidate := range bisects {
		err := gateway.Bisect(candidate[0], candidate[1], candidate[2])
		if err != nil {
			return fmt.Errorf(transmuteSwarmFailFormat, err)
		}
	}

	stitches := swarms[mergeKey]
	for _, candidate := range stitches {
		err := gateway.Stitch(candidate[0], candidate[1], candidate[2])
		if err != nil {
			return fmt.Errorf(transmuteSwarmFailFormat, err)
		}
	}
	return nil
}

/*Edits the swarm map according to the 'candidates' and returns all
successful split/merges that transmuteSwarms can then use to actually
stitch/bisect the p2p swarms*/
func transmuteSwarmMap(swarmMap SwarmMap, candidates []Candidate) (map[int][][]string, error) {
	resultSwarms := map[int][][]string{
		bisectKey: make([][]string, 0),
		mergeKey:  make([][]string, 0),
	}
	for _, candidate := range candidates {
		swarmIDs := candidate.GetSwarmIDs()
		if candidate.IsSplit() {
			idOne, idTwo, err := splitSwarmMap(swarmMap, candidate)
			if err != nil {
				return resultSwarms, err
			}
			resultSwarms[bisectKey] = append(resultSwarms[bisectKey], []string{swarmIDs[0], idOne, idTwo})
		} else {
			id, err := mergeSwarmMap(swarmMap, candidate)
			if err != nil {
				return resultSwarms, err
			}
			resultSwarms[mergeKey] = append(resultSwarms[mergeKey], []string{swarmIDs[0], swarmIDs[1], id})
		}
	}
	return resultSwarms, nil
}

func splitSwarmMap(swarmMap SwarmMap, candidate Candidate) (string, string, error) {
	swarmID := candidate.GetSwarmIDs()[0]
	dspaceOne, dspaceTwo, err := placeDataspaces(swarmMap, swarmID,
		candidate.GetPlacementOne(), candidate.GetPlacementTwo())
	if err != nil {
		return "", "", fmt.Errorf(splitSwarmFailFormat, err)
	}
	swarmIDOne, err := swarmMap.AddSwarm(dspaceOne)
	if err != nil {
		return "", "", fmt.Errorf(splitSwarmFailFormat, err)
	}
	swarmIDTwo, err := swarmMap.AddSwarm(dspaceTwo)
	if err != nil {
		swarmMap.RemoveSwarm(swarmIDOne)
		return "", "", fmt.Errorf(splitSwarmFailFormat, err)
	}

	err = swarmMap.RemoveSwarm(swarmID)
	if err != nil {
		return "", "", fmt.Errorf(splitSwarmFailFormat, err)
	}
	return swarmIDOne, swarmIDTwo, nil
}

func placeDataspaces(swarmMap SwarmMap, swarmID string, placementOne map[string]bool,
	placementTwo map[string]bool) ([]string, []string, error) {
	dataspaces, err := swarmMap.GetDataspaces(swarmID)
	if err != nil {
		return nil, nil, fmt.Errorf(splitSwarmFailFormat, err)
	}

	swarmOne := make([]string, 0)
	swarmTwo := make([]string, 0)
	for _, dataspace := range dataspaces {
		if _, ok := placementOne[dataspace]; ok {
			swarmOne = append(swarmOne, dataspace)
		} else if _, ok := placementTwo[dataspace]; ok {
			swarmTwo = append(swarmTwo, dataspace)
		} else if len(swarmOne) < len(swarmTwo) {
			swarmOne = append(swarmOne, dataspace)
		} else {
			swarmTwo = append(swarmTwo, dataspace)
		}
	}
	return swarmOne, swarmTwo, nil
}

func mergeSwarmMap(swarmMap SwarmMap, candidate Candidate) (string, error) {
	swarms := candidate.GetSwarmIDs()
	swarmIDOne, swarmIDTwo := swarms[0], swarms[1]

	dataspacesOne, err := swarmMap.GetDataspaces(swarmIDOne)
	if err != nil {
		return "", fmt.Errorf(mergeSwarmFailFormat, err)
	}
	dataspacesTwo, err := swarmMap.GetDataspaces(swarmIDTwo)
	if err != nil {
		return "", fmt.Errorf(mergeSwarmFailFormat, err)
	}

	dataspaces := append(dataspacesOne, dataspacesTwo...)
	newSwarmID, err := swarmMap.AddSwarm(dataspaces)
	if err != nil {
		return "", fmt.Errorf(mergeSwarmFailFormat, err)
	}
	swarmMap.RemoveSwarm(swarmIDOne)
	swarmMap.RemoveSwarm(swarmIDTwo)
	return newSwarmID, nil
}
