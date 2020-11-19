package mapper

import (
	"fmt"
	"math"
	"sync"
)

//NewSwarmID takes a seed and creates a new swarm id
var NewSwarmID func(seed int) string = nil

/*SwarmMap is a BiMap of Swarm IDs to managers/dataspaces and
in reverse from dataspaces to Swarm IDs*/
type SwarmMap struct {
	mapMutex              *sync.Mutex
	managerMap            map[string]*swarmMapPair
	inverseMap            map[string]string
	generator             SwarmManagerGenerator
	idCounter             int
	minDataspaceSwarmSize int
	minDataspacesSwarmID  string
}

//New creates a new instance of SwarmMap
func New(generator SwarmManagerGenerator) *SwarmMap {
	return &SwarmMap{
		mapMutex:              &sync.Mutex{},
		managerMap:            make(map[string]*swarmMapPair),
		inverseMap:            make(map[string]string),
		generator:             generator,
		idCounter:             0,
		minDataspaceSwarmSize: math.MaxInt32,
		minDataspacesSwarmID:  "",
	}
}

/*RemoveSwarm removes 'swarmID' from the SwarmMap and cleans up
resources*/
func (sm *SwarmMap) RemoveSwarm(swarmID string) error {
	sm.mapMutex.Lock()
	defer sm.mapMutex.Unlock()

	pair, ok := sm.managerMap[swarmID]
	if !ok {
		return fmt.Errorf("No swarm with ID %s in SwarmMap.RemoveSwarm()", swarmID)
	}
	err := pair.Manager.Close()
	if err != nil {
		return fmt.Errorf("Failed to close swarm %s in SwarmMap.RemoveSwarm(): %v", swarmID, err)
	}

	for _, dspace := range pair.Dataspaces {
		delete(sm.inverseMap, dspace)
	}
	delete(sm.managerMap, swarmID)

	if swarmID == sm.minDataspacesSwarmID {
		sm.minDataspacesSwarmID = ""
		sm.minDataspaceSwarmSize = math.MaxInt32
	}
	return nil
}

/*AddSwarm creates a new swarm with with the 'dataspaces' and returns
the new swarms ID*/
func (sm *SwarmMap) AddSwarm(dataspaces []string) (string, error) {
	sm.mapMutex.Lock()
	defer sm.mapMutex.Unlock()

	newSwarmID := NewSwarmID(sm.idCounter)
	manager, err := sm.generator.New(newSwarmID)
	if err != nil {
		return "", fmt.Errorf("Failed to create new SwarmManager in SwarmMap.AddSwarm(): %v", err)
	}
	sm.idCounter++

	pair := swarmMapPair{Manager: manager.(SwarmManager), Dataspaces: dataspaces}
	sm.managerMap[newSwarmID] = &pair
	for _, dspace := range dataspaces {
		sm.inverseMap[dspace] = newSwarmID
	}

	updateMinDataspacesSwarmInfo(sm, newSwarmID, len(dataspaces))
	return newSwarmID, nil
}

//GetDataspaces returns the dataspaces associated with 'swarmID'
func (sm *SwarmMap) GetDataspaces(swarmID string) ([]string, error) {
	sm.mapMutex.Lock()
	defer sm.mapMutex.Unlock()

	if pair, ok := sm.managerMap[swarmID]; ok {
		return pair.Dataspaces, nil
	}
	return nil, fmt.Errorf("No swarm with ID %s in SwarmMap.GetDataspaces()", swarmID)
}

//GetSwarmID returns the swarm ID associated with a dataspace
func (sm *SwarmMap) GetSwarmID(dspaceID string) (string, error) {
	sm.mapMutex.Lock()
	defer sm.mapMutex.Unlock()

	if swarmID, ok := sm.inverseMap[dspaceID]; ok {
		return swarmID, nil
	}
	return "", fmt.Errorf("No swarm associated with dataspace %s in SwarmMap.GetSwarmID()", dspaceID)
}

//GetSwarmManager returns the swarm manager object associated with 'swarmID'
func (sm *SwarmMap) GetSwarmManager(swarmID string) (interface{}, error) {
	sm.mapMutex.Lock()
	defer sm.mapMutex.Unlock()

	if pair, ok := sm.managerMap[swarmID]; ok {
		return pair.Manager, nil
	}
	return nil, fmt.Errorf("No swarm with ID %s in SwarmMap.GetSwarmManager()", swarmID)
}

/*GetMinDataspaceSwarm returns the swarm ID of the swarm with the least
associated dataspaces*/
func (sm *SwarmMap) GetMinDataspaceSwarm() (string, error) {
	sm.mapMutex.Lock()
	defer sm.mapMutex.Unlock()

	if len(sm.managerMap) == 0 {
		return "", fmt.Errorf("No swarm with minimum dataspaces in SwarmMap.GetMinDataspaceSwarm()")
	}
	return calculateMinDataspaceSwarm(sm), nil
}

//AddDataspace adds 'dataspace' to 'swarmID'
func (sm *SwarmMap) AddDataspace(swarmID string, dataspace string) error {
	sm.mapMutex.Lock()
	defer sm.mapMutex.Unlock()

	if pair, ok := sm.managerMap[swarmID]; ok {
		pair.Dataspaces = append(pair.Dataspaces, dataspace)
		return nil
	}
	return fmt.Errorf("No swarm with ID %s in SwarmMap.AddDataspace()", swarmID)
}

//RemoveDataspace removes 'dataspace' from 'swarmID'
func (sm *SwarmMap) RemoveDataspace(swarmID string, dataspace string) error {
	sm.mapMutex.Lock()
	defer sm.mapMutex.Unlock()

	if pair, ok := sm.managerMap[swarmID]; ok {
		dspacesLen := len(pair.Dataspaces)
		for i := 0; i < dspacesLen; i++ {
			if pair.Dataspaces[i] == dataspace {
				pair.Dataspaces[i] = pair.Dataspaces[dspacesLen-1]
				pair.Dataspaces = pair.Dataspaces[:dspacesLen-1]
				updateMinDataspacesSwarmInfo(sm, swarmID, len(pair.Dataspaces))
				return nil
			}
		}
		return fmt.Errorf("No dataspace with ID %s associated with swarm %s in"+
			"SwarmMap.RemoveDataspace()", dataspace, swarmID)
	}
	return fmt.Errorf("No swarm with ID %s in SwarmMap.AddDataspace()", swarmID)
}
