package transmuter

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func TestSwarmTransmuter(t *testing.T) {
	totalSwarms := 5
	endpointsPerSwarm := 5
	smap := &TestSwarmMap{managers: make(map[string]SwarmManager)}
	analyzer := &TestSwarmAnalyzer{smap: smap}
	for i := 0; i < totalSwarms; i++ {
		endpoints := make([]string, endpointsPerSwarm)
		for j := 0; j < endpointsPerSwarm; j++ {
			endpoints[j] = "/endpoint/" + strconv.Itoa(i*endpointsPerSwarm+j)
		}
		swarmID := "/dataspace/" + strconv.Itoa(i)
		manager := &TestSwarmManager{endpoints: endpoints}
		smap.managers[swarmID] = manager
	}
	printSwarmSizes(smap.managers)

	PollPeriod = time.Second
	transmuter := New(smap, analyzer)

	totalConnections := 50
	for i := 0; i < totalConnections; i++ {
		fc := &FakeConn{id: "/endpoint/" + strconv.Itoa(totalSwarms*endpointsPerSwarm+i)}
		dspace := "/dataspace/" + strconv.Itoa(rand.Intn(totalSwarms))
		transmuter.ProcessConnection(dspace, 0, fc)
	}
	printSwarmSizes(smap.managers)
	time.Sleep(time.Second * 5)
	printSwarmSizes(smap.managers)
}

func printSwarmSizes(m map[string]SwarmManager) {
	fmt.Printf("\n---------------------------------\n")
	for id, manager := range m {
		fmt.Printf("Swarm: %s\n", id)
		endpoints := manager.GetEndpoints()
		for _, endpoint := range endpoints {
			fmt.Printf("\t%s\n", endpoint)
		}
	}
}

type TestSwarmMap struct {
	managers map[string]SwarmManager
}

func (tm *TestSwarmMap) GetSwarm(id string) (interface{}, error) {
	if manager, ok := tm.managers[id]; ok {
		return manager, nil
	}
	return nil, fmt.Errorf("No Swarm with id %s", id)
}

type TestCandidate struct {
	transferer string
	transferee string
	size       int
}

func (tc *TestCandidate) GetTransfererID() string { return tc.transferer }
func (tc *TestCandidate) GetTransfereeID() string { return tc.transferee }
func (tc *TestCandidate) GetTransferSize() int    { return tc.size }

type TestSwarmAnalyzer struct {
	smap *TestSwarmMap
}

func (ta *TestSwarmAnalyzer) GetMostNeedy() (string, error) {
	iterLen := rand.Intn(len(ta.smap.managers))
	i := 0
	var id string
	for id, _ = range ta.smap.managers {
		if i == iterLen {
			break
		}
		i++
	}
	return id, nil
}

func (ta *TestSwarmAnalyzer) CalculateCandidates() ([]Candidate, error) {
	totalPairings := rand.Intn(len(ta.smap.managers) / 2)
	candidates := make([]Candidate, totalPairings)
	slice := mapToSlice(ta.smap.managers)
	for i := 0; i < totalPairings; i++ {
		id1, id2 := slice[i], slice[len(slice)-i-1]
		size := rand.Intn(10)
		candidates[i] = &TestCandidate{transferer: id1, transferee: id2, size: size}
	}
	return candidates, nil
}

func mapToSlice(m map[string]SwarmManager) []string {
	s := make([]string, 0, len(m))
	for key, _ := range m {
		s = append(s, key)
	}
	return s
}

type TestSwarmManager struct {
	endpoints []string
}

func (sm *TestSwarmManager) SetID(string) {}
func (sm *TestSwarmManager) AddEndpoint(i interface{}) error {
	c := i.(*FakeConn)
	return sm.TakeEndpoint(c.id)
}
func (sm *TestSwarmManager) RemoveEndpoint(i interface{}) error {
	c := i.(*FakeConn)
	return sm.DropEndpoint(c.id)
}
func (sm *TestSwarmManager) TakeEndpoint(s string) error {
	sm.endpoints = append(sm.endpoints, s)

	return nil
}
func (sm *TestSwarmManager) DropEndpoint(s string) error {
	for i := 0; i < len(sm.endpoints); i++ {
		if sm.endpoints[i] == s {
			sm.endpoints = append(sm.endpoints[:i], sm.endpoints[i+1:]...)
			return nil
		}
	}
	return nil
}
func (sm *TestSwarmManager) GetEndpoints() []string {
	return sm.endpoints
}
func (sm *TestSwarmManager) Close() error {
	return nil
}

type FakeConn struct{ id string }

func (fc *FakeConn) Read([]byte) (int, error)  { return 0, nil }
func (fc *FakeConn) Write([]byte) (int, error) { return 0, nil }
func (fc *FakeConn) Close() error              { return nil }
