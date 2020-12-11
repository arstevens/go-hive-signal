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
		manager := New("/swarm/"+strconv.Itoa(i), gateway, negotiate, tracker)
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
	}

	fmt.Printf("[RUNNING TRANSMUTATION TESTS]\n")
	fmt.Printf("Merges\n---------\n")
	for i := 0; i < totalSwarms/2; i++ {
		s1 := swarms[i]
		s2 := swarms[totalSwarms/2+i]
		fmt.Printf("\tOriginal Sizes: %d %d\n", s1.gateway.GetTotalEndpoints(), s2.gateway.GetTotalEndpoints())
		err := s1.Stitch(s2)
		if err != nil {
			panic(err)
		}
		fmt.Printf("\tMerge Size: %d\n", s1.gateway.GetTotalEndpoints())
	}
	swarms = swarms[:totalSwarms/2]
	fmt.Printf("Bisects\n----------\n")
	for i := 0; i < totalSwarms/2; i++ {
		s1 := swarms[i]
		fmt.Printf("\tOriginal Size: %d\n", s1.gateway.GetTotalEndpoints())
		s2, err := s1.Bisect()
		if err != nil {
			panic(err)
		}
		s := s2.(*SwarmManager)
		fmt.Printf("\tSplit Sizes: %d %d\n", s1.gateway.GetTotalEndpoints(), s.gateway.GetTotalEndpoints())
	}

}

type testSwarmTracker struct {
	m map[string]int
}

func (st *testSwarmTracker) SetSize(s string, i int) {
	st.m[s] = i
}

type testSwarmGateway struct {
	conn           *FakeConn
	totalEndpoints int
}

func (sg *testSwarmGateway) GetEndpoint() (Conn, error) {
	return sg.conn, nil
}
func (sg *testSwarmGateway) AddEndpoint(Conn) error {
	sg.totalEndpoints++
	return nil
}
func (sg *testSwarmGateway) RetireEndpoint(Conn) error {
	if sg.totalEndpoints == 0 {
		return fmt.Errorf("No endpoint to retire")
	}
	sg.totalEndpoints--
	return nil
}
func (sg *testSwarmGateway) EvenlySplit() (SwarmGateway, error) {
	sg.totalEndpoints /= 2
	return &testSwarmGateway{conn: sg.conn, totalEndpoints: sg.totalEndpoints}, nil
}
func (sg *testSwarmGateway) GetTotalEndpoints() int {
	return sg.totalEndpoints
}
func (sg *testSwarmGateway) Merge(g SwarmGateway) error {
	gway := g.(*testSwarmGateway)
	sg.totalEndpoints += gway.totalEndpoints
	gway.Close()
	return nil
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
