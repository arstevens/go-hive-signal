package manager

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
)

func TestManager(t *testing.T) {
	tracker := &testSwarmTracker{m: make(map[string]int)}
	totalSwarms := 10
	swarms := make([]*SwarmManager, totalSwarms)
	fmt.Printf("[GENERATING %d SWARMS]\n", totalSwarms)
	for i := 0; i < totalSwarms; i++ {
		gateway := &testSwarmGateway{conn: &FakeConn{}, totalEndpoints: rand.Intn(100)}
		manager := New("/dataspace/"+strconv.Itoa(i), gateway, negotiate, tracker)
		swarms[i] = manager
	}

	fmt.Printf("[RUNNING MANAGER TESTS]\n")
	for i := 0; i < totalSwarms; i++ {
		fmt.Printf("\tID: %s\n", swarms[i].GetID())
		conn := &FakeConn{}
		err := swarms[i].AttemptToPair(conn)
		if err != nil {
			t.Fatalf("Failed to run pair: %v\n", err)
		}
		err = swarms[i].AddEndpoint(&FakeConn{})
		if err != nil {
			panic(err)
		}
		err = swarms[i].RemoveEndpoint(&FakeConn{})
		if err != nil {
			panic(err)
		}
		err = swarms[i].TakeEndpoint("")
		if err != nil {
			panic(err)
		}
		err = swarms[i].DropEndpoint("")
		if err != nil {
			panic(err)
		}
	}
}

type testSwarmTracker struct {
	m map[string]int
}

func (st *testSwarmTracker) SetSize(s string, i int) {
	st.m[s] = i
}

func (st *testSwarmTracker) Delete(s string) {
	delete(st.m, s)
}

func (st *testSwarmTracker) AddDebriefDatapoint(string, interface{}) {}

type testSwarmGateway struct {
	conn           *FakeConn
	totalEndpoints int
}

func (sg *testSwarmGateway) GetEndpoint() (Conn, interface{}, error) {
	return sg.conn, rand.Intn(95) + 5, nil
}
func (sg *testSwarmGateway) PushEndpoint(string) error {
	sg.totalEndpoints++
	return nil
}
func (sg *testSwarmGateway) RemoveEndpoint(string) error {
	if sg.totalEndpoints == 0 {
		return fmt.Errorf("No endpoint to retire")
	}
	sg.totalEndpoints--
	return nil
}
func (sg *testSwarmGateway) GetTotalEndpoints() int {
	return sg.totalEndpoints
}
func (sg *testSwarmGateway) GetEndpointAddrs() []string {
	addrs := make([]string, sg.totalEndpoints)
	for i := 0; i < sg.totalEndpoints; i++ {
		addrs[i] = ""
	}
	return addrs
}
func (sg *testSwarmGateway) Close() error {
	sg.totalEndpoints = 0
	return nil
}

func negotiate(Conn, Conn) error {
	return nil
}

type FakeConn struct{}

func (fc *FakeConn) Read([]byte) (int, error)  { return 0, nil }
func (fc *FakeConn) Write([]byte) (int, error) { return 0, nil }
func (fc *FakeConn) Close() error              { return nil }
func (fc *FakeConn) GetAddress() string        { return "" }
func (fc *FakeConn) IsClosed() bool            { return false }
