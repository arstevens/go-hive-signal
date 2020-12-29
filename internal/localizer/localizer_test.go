package localizer

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestLocalizer(t *testing.T) {
	totalRequests := 50
	totalDataspaces := 5
	requests := make([]LocalizeRequestTest, totalRequests)
	for i := 0; i < totalRequests; i++ {
		requests[i] = LocalizeRequestTest{
			d: fmt.Sprintf("/dataspace/%d", rand.Intn(totalDataspaces)),
		}
	}

	smap := SwarmMapTest{smap: make(map[string]SwarmManager)}
	for i := 0; i < totalDataspaces; i++ {
		dataspace := fmt.Sprintf("/dataspace/%d", i)
		smap.smap[dataspace] = &SwarmManagerTest{id: dataspace}
	}
	ftrack := FrequencyTrackerTest{fmap: make(map[string]int)}
	fconn := FakeConn{}

	queueSize := 3
	rlocalizer := New(queueSize, &smap, &ftrack)
	for _, request := range requests {
		rlocalizer.AddJob(&LocalizeRequestTest{
			d: request.GetDataspace(),
		}, &fconn)
	}
	time.Sleep(time.Second)
}

type FakeConn struct{}

func (fc *FakeConn) Read([]byte) (int, error)  { return 0, nil }
func (fc *FakeConn) Write([]byte) (int, error) { return 0, nil }
func (fc *FakeConn) Close() error              { return nil }

type SwarmMapTest struct {
	smap map[string]SwarmManager
}

func (st *SwarmMapTest) GetSwarm(s string) (interface{}, error) {
	for _, manager := range st.smap {
		if manager.GetID() == s {
			return manager, nil
		}
	}
	return nil, fmt.Errorf("No swarm with name %s", s)
}

type SwarmManagerTest struct {
	id string
}

func (sm *SwarmManagerTest) AttemptToPair(conn interface{}) error {
	fmt.Printf("%s: Attempting to pair\n", sm.id)
	return nil
}

func (sm *SwarmManagerTest) GetID() string {
	return sm.id
}

type FrequencyTrackerTest struct {
	fmap map[string]int
}

func (ft *FrequencyTrackerTest) IncrementFrequency(dataspace string, s string) {
	ft.fmap[s]++
	fmt.Printf("Swarm %s incremented to %d\n", s, ft.fmap[s])
}

type LocalizeRequestTest struct {
	d string
}

func (lt *LocalizeRequestTest) GetDataspace() string { return lt.d }
