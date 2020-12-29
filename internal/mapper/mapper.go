package mapper

import (
	"fmt"
	"io"
	"sync"
)

/*SwarmManagerGenerator describes an object that can create
new SwarmManagers given an id*/
type SwarmManagerGenerator interface {
	New(id string) interface{}
}

/*SwarmMap holds a thread-safe mapping of dataspaces
to swarm managers*/
type SwarmMap struct {
	mapMutex   *sync.Mutex
	managerMap map[string]interface{}
	generator  SwarmManagerGenerator
}

//New creates a new instance of SwarmMap
func New(generator SwarmManagerGenerator) *SwarmMap {
	return &SwarmMap{
		mapMutex:   &sync.Mutex{},
		managerMap: make(map[string]interface{}),
		generator:  generator,
	}
}

/*RemoveSwarm removes dataspace from the SwarmMap and cleans up
resources*/
func (sm *SwarmMap) RemoveSwarm(dataspace string) error {
	sm.mapMutex.Lock()
	defer sm.mapMutex.Unlock()

	if manager, ok := sm.managerMap[dataspace]; ok {
		closer := manager.(io.Closer)
		closer.Close()
		delete(sm.managerMap, dataspace)
		return nil
	}
	return fmt.Errorf("No swarm associated with dataspace %s in SwarmMap.RemoveSwarm()", dataspace)
}

/*AddSwarm creates a new swarm associated with the dataspace*/
func (sm *SwarmMap) AddSwarm(dataspace string) error {
	manager := sm.generator.New(dataspace)

	sm.mapMutex.Lock()
	sm.managerMap[dataspace] = manager
	sm.mapMutex.Unlock()
	return nil
}

//GetSwarm returns the swarm manager object associated with the dataspace
func (sm *SwarmMap) GetSwarm(dataspace string) (interface{}, error) {
	sm.mapMutex.Lock()
	defer sm.mapMutex.Unlock()

	if manager, ok := sm.managerMap[dataspace]; ok {
		return manager, nil
	}
	return nil, fmt.Errorf("No swarm associated with dataspace %s in SwarmMap.GetSwarmManager()", dataspace)
}
