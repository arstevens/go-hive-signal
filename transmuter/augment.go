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
func transmuteSwarmMap(swarmMap SwarmMap, candidates [][]string) (map[int][][]string, error) {
	resultSwarms := map[int][][]string{
		bisectKey: make([][]string, 0),
		mergeKey:  make([][]string, 0),
	}
	for _, candidate := range candidates {
		if len(candidate) == 1 {
			idOne, idTwo, err := splitSwarmMap(swarmMap, candidate[0])
			if err != nil {
				return resultSwarms, err
			}
			resultSwarms[bisectKey] = append(resultSwarms[bisectKey], []string{candidate[0], idOne, idTwo})
		} else if len(candidate) == 2 {
			id, err := mergeSwarmMap(swarmMap, candidate[0], candidate[1])
			if err != nil {
				return resultSwarms, err
			}
			resultSwarms[mergeKey] = append(resultSwarms[mergeKey], []string{candidate[0], candidate[1], id})
		}
	}
	return resultSwarms, nil
}

func splitSwarmMap(swarmMap SwarmMap, swarmID string) (string, string, error) {
	dataspaces, err := swarmMap.GetDataspaces(swarmID)
	if err != nil {
		return "", "", fmt.Errorf(splitSwarmFailFormat, err)
	}
	dspaceOne, dspaceTwo := splitStringSlice(dataspaces)
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

func splitStringSlice(slice []string) ([]string, []string) {
	midpoint := len(slice) / 2
	s1 := slice[:midpoint]
	s2 := slice[midpoint:]
	return s1, s2
}

func mergeSwarmMap(swarmMap SwarmMap, swarmIDOne string, swarmIDTwo string) (string, error) {
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
