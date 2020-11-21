package manager

import (
	"fmt"
	"strconv"
	"testing"
)

func TestManager(t *testing.T) {
	gateway := testSwarmGateway{&FakeConn{}}
	fmt.Printf("[CREATING GENERATOR]\n")
	generator := NewGenerator(&gateway, negotiate)
	if generator == nil {
		t.Fatal("Failed to create generator")
	}

	totalSwarms := 10
	swarms := make([]*SwarmManager, totalSwarms)
	fmt.Printf("[GENERATING %d SWARMS]\n", totalSwarms)
	for i := 0; i < totalSwarms; i++ {
		manager, err := generator.New("/swarm/" + strconv.Itoa(i))
		if err != nil {
			t.Fatalf("Failed to generate swarms: %v\n", err)
		}
		swarms[i] = manager.(*SwarmManager)
	}

	fmt.Printf("[RUNNING MANAGER TESTS]\n")
	for i := 0; i < totalSwarms; i++ {
		fmt.Printf("\tID: %s\n", swarms[i].GetID())
		conn := &FakeConn{}
		err := swarms[i].AttemptToPair(conn)
		if err != nil {
			t.Fatalf("Failed to run pair: %v\n", err)
		}
		swarms[i].Close()
		err = swarms[i].AttemptToPair(conn)
		if err == nil {
			t.Fatalf("Failed to fail when closed: %v\n", err)
		}
	}
}

type testSwarmGateway struct {
	conn *FakeConn
}

func (sg *testSwarmGateway) GetEndpoint(string) (interface{}, error) {
	return sg.conn, nil
}

func negotiate(Conn, Conn) error {
	return nil
}

type FakeConn struct{}

func (fc *FakeConn) Read([]byte) (int, error)  { return 0, nil }
func (fc *FakeConn) Write([]byte) (int, error) { return 0, nil }
func (fc *FakeConn) Close() error              { return nil }
