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

func pollForTransmutation(swarmMap SwarmMap, analyzer SwarmAnalyzer) {
	for {
		time.Sleep(PollPeriod)
		candidates, err := analyzer.CalculateCandidates()
		if err != nil {
			return
		}

		if len(candidates) > 0 {
			newSwarms, oldSwarms, err := transmuteSwarms(swarmMap, candidates)
			if err != nil {
				log.Println(err)
			}
			err = transmuteSwarmMap(swarmMap, oldSwarms, newSwarms)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

type transmuteReturnPair struct {
	manager    SwarmManager
	dataspaces []string
}

func transmuteSwarms(swarmMap SwarmMap, candidates []Candidate) ([]*transmuteReturnPair, []string, error) {
	newManagers := make([]*transmuteReturnPair, 0)
	oldSwarms := make([]string, 0)
	for _, candidate := range candidates {
		if candidate.IsSplit() {
			pairOne, pairTwo, err := splitTransmutation(swarmMap, candidate)
			if err != nil {
				return nil, nil, err
			}
			newManagers = append(newManagers, pairOne, pairTwo)
			oldSwarms = append(oldSwarms, candidate.GetSwarmIDs()[0])
		} else {
			swarmIDs := candidate.GetSwarmIDs()
			newPair, err := mergeTransmutation(swarmMap, swarmIDs)
			if err != nil {
				return nil, nil, err
			}
			newManagers = append(newManagers, newPair)
			oldSwarms = append(oldSwarms, swarmIDs...)
		}
	}
	return newManagers, oldSwarms, nil
}

func mergeTransmutation(swarmMap SwarmMap, swarmIDs []string) (*transmuteReturnPair, error) {
	managerOne, err := swarmMap.GetSwarmByID(swarmIDs[0])
	if err != nil {
		return nil, fmt.Errorf(transmuteSwarmFailFormat, err)
	}
	managerTwo, err := swarmMap.GetSwarmByID(swarmIDs[1])
	if err != nil {
		return nil, fmt.Errorf(transmuteSwarmFailFormat, err)
	}
	err = managerOne.Stitch(managerTwo)
	if err != nil {
		return nil, fmt.Errorf(transmuteSwarmFailFormat, err)
	}
	managerTwo.Close()

	dataspaces, err := consolidateDataspaces(swarmMap, swarmIDs[0], swarmIDs[1])
	if err != nil {
		return nil, fmt.Errorf(transmuteSwarmFailFormat, err)
	}
	newPair := transmuteReturnPair{manager: managerOne, dataspaces: dataspaces}
	return &newPair, nil
}

func splitTransmutation(swarmMap SwarmMap, candidate Candidate) (*transmuteReturnPair, *transmuteReturnPair, error) {
	swarmID := candidate.GetSwarmIDs()[0]
	manager, err := swarmMap.GetSwarmByID(swarmID)
	if err != nil {
		return nil, nil, fmt.Errorf(transmuteSwarmFailFormat, err)
	}
	newManager, err := manager.Bisect()
	if err != nil {
		return nil, nil, fmt.Errorf(transmuteSwarmFailFormat, err)
	}
	dspacesOne, dspacesTwo, err := placeDataspaces(swarmMap, swarmID, candidate.GetPlacementOne(), candidate.GetPlacementTwo())
	if err != nil {
		return nil, nil, fmt.Errorf(transmuteSwarmFailFormat, err)
	}
	pairOne := transmuteReturnPair{manager: manager, dataspaces: dspacesOne}
	pairTwo := transmuteReturnPair{manager: newManager, dataspaces: dspacesTwo}
	return &pairOne, &pairTwo, nil
}

/*Edits the swarm map according to the 'candidates' and returns all
successful split/merges that transmuteSwarms can then use to actually
stitch/bisect the p2p swarms*/
func transmuteSwarmMap(swarmMap SwarmMap, oldSwarms []string, newSwarms []*transmuteReturnPair) error {
	var err error
	for _, swarmID := range oldSwarms {
		err = swarmMap.RemoveSwarm(swarmID)
		if err != nil {
			return fmt.Errorf(transmuteSwarmFailFormat, err)
		}
	}

	for _, pair := range newSwarms {
		newID, err := swarmMap.AddSwarm(pair.manager, pair.dataspaces)
		if err != nil {
			return fmt.Errorf(transmuteSwarmFailFormat, err)
		}
		pair.manager.SetID(newID)
	}
	return nil
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

func consolidateDataspaces(swarmMap SwarmMap, swarmIDOne string, swarmIDTwo string) ([]string, error) {
	dataspacesOne, err := swarmMap.GetDataspaces(swarmIDOne)
	if err != nil {
		return nil, fmt.Errorf(mergeSwarmFailFormat, err)
	}
	dataspacesTwo, err := swarmMap.GetDataspaces(swarmIDTwo)
	if err != nil {
		return nil, fmt.Errorf(mergeSwarmFailFormat, err)
	}

	dataspaces := append(dataspacesOne, dataspacesTwo...)
	return dataspaces, nil
}
