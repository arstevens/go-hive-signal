package mapper

import (
	"fmt"
	"strconv"
	"testing"
)

func TestMapper(t *testing.T) {
	generator := &testGenerator{}
	swarmMapper := New(generator)

	// Test Swarm Add
	totalSwarms := 100
	dspaces := make([]string, 0)
	for i := 0; i < totalSwarms; i++ {
		dataspace := "/dataspace/" + strconv.Itoa(i)
		err := swarmMapper.AddSwarm(dataspace)
		if err != nil {
			t.Fatal(err)
		}
		dspaces = append(dspaces, dataspace)
	}
	fmt.Printf("Added Swarms (%d-%d)\n", 0, totalSwarms-1)

	// Test GetSwarm
	_, err := swarmMapper.GetSwarm(dspaces[0])
	if err != nil {
		panic(err)
	}
	fmt.Printf("Successfully retrieved manager for (%s)\n", dspaces[0])

	// Test RemoveSwarm
	err = swarmMapper.RemoveSwarm(dspaces[0])
	if err != nil {
		panic(err)
	}
	fmt.Printf("Removed swarm (%s)\n", dspaces[0])
	_, err = swarmMapper.GetSwarm(dspaces[0])
	if err == nil {
		t.Fatalf("Swarm should have been deleted")
	}

}

type TestSwarmManager struct{}

func (tm *TestSwarmManager) Close() error {
	return nil
}

type testGenerator struct{}

func (tg *testGenerator) New(string) interface{} { return &TestSwarmManager{} }
