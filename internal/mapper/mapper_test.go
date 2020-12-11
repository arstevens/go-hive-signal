package mapper

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/arstevens/go-hive-signal/internal/transmuter"
	"github.com/arstevens/go-request/handle"
)

func TestMapper(t *testing.T) {
	NewSwarmID = func(i int) string { return "/swarm/" + strconv.Itoa(i) }
	swarmMapper := New()

	// Test Swarm Add
	totalSwarms := 100
	totalDspacesPer := 60
	allDspaces := make([]string, 0)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < totalSwarms; i++ {
		totDspaces := rand.Intn(totalDspacesPer-1) + 1
		dspaces := make([]string, totDspaces)
		for j := 0; j < totDspaces; j++ {
			dspaces[j] = "/dataspace/" + strconv.Itoa(i*totalDspacesPer+j)
			allDspaces = append(allDspaces, dspaces[j])
		}
		_, err := swarmMapper.AddSwarm(&TestSwarmManager{}, dspaces)
		if err != nil {
			panic(err)
		}
	}
	fmt.Printf("Added Swarms (%d-%d)\n", 0, totalSwarms-1)

	// Test GetMin and GetDataspaces
	minID, err := swarmMapper.GetMinDataspaceSwarm()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Min Swarm: %s with size %d\n", minID, len(swarmMapper.managerMap[minID].Dataspaces))
	dspaces, _ := swarmMapper.GetDataspaces(minID)
	fmt.Printf("\t(Dataspaces)%v\n", dspaces)

	// Test GetSwarmID and GetSwarmManager
	sid, err := swarmMapper.GetSwarmID(dspaces[0])
	if err != nil {
		panic(err)
	}
	fmt.Printf("Swarm associated with (%s) is (%s)\n", dspaces[0], sid)
	_, err = swarmMapper.GetSwarmManager(sid)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Successfully retrieved manager for (%s)\n", sid)

	// Test RemoveDataspace
	err = swarmMapper.RemoveDataspace(sid, dspaces[0])
	if err != nil {
		panic(err)
	}
	fmt.Printf("Removed dataspace (%s) from (%s)\n", dspaces[0], sid)
	dspaces, _ = swarmMapper.GetDataspaces(sid)
	fmt.Printf("\t(Dataspaces)%v\n", dspaces)

	// Test AddDataspace
	newDspace := "/dataspace/ADDITION"
	err = swarmMapper.AddDataspace(sid, newDspace)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Added dataspace (%s) to (%s)\n", newDspace, sid)
	dspaces, _ = swarmMapper.GetDataspaces(sid)
	fmt.Printf("\t(Dataspaces)%v\n", dspaces)

	// Test RemoveSwarm
	err = swarmMapper.RemoveSwarm(sid)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Removed swarm (%s)\n", sid)
	_, err = swarmMapper.GetDataspaces(sid)
	if err == nil {
		t.Fatalf("Swarm should have been deleted")
	}

}

type TestSwarmManager struct{}

func (tm *TestSwarmManager) Close() error {
	return nil
}
func (tm *TestSwarmManager) AddEndpoint(handle.Conn) error    { return nil }
func (tm *TestSwarmManager) RemoveEndpoint(handle.Conn) error { return nil }
func (tm *TestSwarmManager) Bisect() (transmuter.SwarmManager, error) {
	return &TestSwarmManager{}, nil
}
func (tm *TestSwarmManager) Stitch(man transmuter.SwarmManager) error {
	man.Close()
	return nil
}
