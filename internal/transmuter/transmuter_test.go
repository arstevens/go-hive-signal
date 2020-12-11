package transmuter

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/arstevens/go-request/handle"
)

func TestTransmuterSplitMerge(t *testing.T) {
	fmt.Println("\n----------STARTING SPLIT/MERGE TEST-------------")
	totalStartingSwarms := 10
	totalDspacesPerSwarm := 5

	swarmMap := TestSwarmMap{swarms: make(map[string][]string), managers: make(map[string]SwarmManager), idCount: totalStartingSwarms}
	dspaceCount := 0
	for i := 0; i < totalStartingSwarms; i++ {
		id := "/swarm/" + strconv.Itoa(i)
		dspaces := make([]string, totalDspacesPerSwarm)
		for j := 0; j < totalDspacesPerSwarm; j++ {
			dspaces[j] = "/dataspace/" + strconv.Itoa(dspaceCount)
			dspaceCount++
		}
		swarmMap.swarms[id] = dspaces
		swarmMap.managers[id] = &TestSwarmManager{}
	}
	fmt.Printf("Pre-transmutation Mapping\n----------------------\n")
	printSwarmMap(swarmMap)
	tracker := TestSizeTracker{sizes: make(map[string]int)}
	analyzer := TestSwarmAnalyzer{smap: &swarmMap}

	PollPeriod = time.Second
	_ = New(&tracker, &swarmMap, &analyzer)
	time.Sleep(time.Second * 5)

	fmt.Printf("Post-transmutation Mapping\n----------------------\n")
	printSwarmMap(swarmMap)
}

func printSwarmMap(smap TestSwarmMap) {
	for id, dspaces := range smap.swarms {
		fmt.Printf("(SwarmID)->%s\n", id)
		for _, dspace := range dspaces {
			fmt.Printf("\t%s\n", dspace)
		}
	}
}

func TestTransmuterAddRemove(t *testing.T) {
	fmt.Println("\n----------STARTING ADD/REMOVE TEST-------------")
	totalRequests := 50
	totalStartingSwarms := 10
	totalDspacesPerSwarm := 5

	swarmMap := TestSwarmMap{swarms: make(map[string][]string), idCount: totalStartingSwarms}
	tracker := TestSizeTracker{sizes: make(map[string]int)}
	dspaceCount := 0
	for i := 0; i < totalStartingSwarms; i++ {
		id := "/swarm/" + strconv.Itoa(i)
		dspaces := make([]string, totalDspacesPerSwarm)
		for j := 0; j < totalDspacesPerSwarm; j++ {
			dspaces[j] = "/dataspace/" + strconv.Itoa(dspaceCount)
			dspaceCount++
		}
		swarmMap.swarms[id] = dspaces
		tracker.sizes[id] = 0
	}
	analyzer := TestSwarmAnalyzer{smap: &swarmMap}

	PollPeriod = time.Hour
	transmuter := New(&tracker, &swarmMap, &analyzer)

	conn := FakeConn{}
	for i := 0; i < totalRequests; i++ {
		var id string
		var code int
		j, counter := 0, rand.Intn(len(swarmMap.swarms))
		for sID, _ := range swarmMap.swarms {
			if j == counter {
				id = sID
				break
			}
			j++
		}

		code = SwarmConnect
		if rand.Intn(100) < 50 {
			code = SwarmDisconnect
		}
		transmuter.ProcessConnection(id, code, &conn)
	}
}

type TestSizeTracker struct {
	sizes map[string]int
}

func (tt *TestSizeTracker) GetSmallest() (string, error) {
	if len(tt.sizes) == 0 {
		return "", fmt.Errorf("No smallest swarms available in TestSizeTracker")
	}

	minSize := math.MaxInt32
	minID := ""
	for id, size := range tt.sizes {
		if size < minSize {
			minSize = size
			minID = id
		}
	}
	fmt.Printf("Returing smallest swarm %s\n", minID)
	return minID, nil
}

type TestSwarmMap struct {
	swarms   map[string][]string
	managers map[string]SwarmManager
	idCount  int
}

func (tm *TestSwarmMap) RemoveSwarm(id string) error {
	if _, ok := tm.swarms[id]; !ok {
		return fmt.Errorf("Swarm %s does not exist in RemoveSwarm", id)
	}
	delete(tm.swarms, id)
	return nil
}

func (tm *TestSwarmMap) AddSwarm(manager SwarmManager, dspaces []string) (string, error) {
	id := "/swarm/" + strconv.Itoa(tm.idCount)
	tm.idCount++

	tm.swarms[id] = dspaces
	tm.managers[id] = manager
	return id, nil
}

func (tm *TestSwarmMap) GetSwarmByID(id string) (SwarmManager, error) {
	if manager, ok := tm.managers[id]; ok {
		return manager, nil
	}
	return nil, fmt.Errorf("No Swarm with id %s", id)
}

func (tm *TestSwarmMap) GetDataspaces(id string) ([]string, error) {
	if _, ok := tm.swarms[id]; !ok {
		return nil, fmt.Errorf("Swarm %s does not exist in GetDataspaces", id)
	}
	return tm.swarms[id], nil
}

type TestCandidate struct {
	split      bool
	ids        []string
	placements []map[string]bool
}

func (tc *TestCandidate) IsSplit() bool                    { return tc.split }
func (tc *TestCandidate) GetSwarmIDs() []string            { return tc.ids }
func (tc *TestCandidate) GetPlacementOne() map[string]bool { return tc.placements[0] }
func (tc *TestCandidate) GetPlacementTwo() map[string]bool { return tc.placements[1] }

func splitStringSlice(s []string) []map[string]bool {
	placements := make([]map[string]bool, 2)
	placements[0] = make(map[string]bool)
	placements[1] = make(map[string]bool)
	for i := 0; i < len(s)/2; i++ {
		placements[0][s[i]] = true
	}
	for i := len(s) / 2; i < len(s); i++ {
		placements[1][s[i]] = true
	}
	return placements
}

type TestSwarmAnalyzer struct {
	smap *TestSwarmMap
}

func (ta *TestSwarmAnalyzer) CalculateCandidates() ([]Candidate, error) {
	candidates := make([]Candidate, 0)
	candidate := make([]string, 0)
	for id, dataspaces := range ta.smap.swarms {
		if len(candidate) == 2 {
			finalCandidate := TestCandidate{split: false, ids: candidate, placements: nil}
			candidates = append(candidates, &finalCandidate)
			candidate = make([]string, 0)
		} else if len(candidate) == 1 && rand.Intn(100) < 50 {
			finalCandidate := TestCandidate{split: true, ids: candidate, placements: splitStringSlice(dataspaces)}
			candidates = append(candidates, &finalCandidate)
			candidate = make([]string, 0)
		}
		candidate = append(candidate, id)
	}
	return candidates, nil
}

type TestSwarmGateway struct{}

func (tg *TestSwarmGateway) AddEndpoint(id string, conn handle.Conn) error {
	fmt.Printf("(Gateway) Adding endpoint to swarm %s\n", id)
	return nil
}
func (tg *TestSwarmGateway) Bisect(id string, newIDOne string, newIDTwo string) error {
	fmt.Printf("(Gateway) Bisecting %s into (%s, %s)\n", id, newIDOne, newIDTwo)
	return nil
}
func (tg *TestSwarmGateway) Stitch(idOne string, idTwo string, newID string) error {
	fmt.Printf("(Gateway) Stitching (%s, %s) into %s\n", idOne, idTwo, newID)
	return nil
}

type TestSwarmManager struct{}

func (sm *TestSwarmManager) AddEndpoint(handle.Conn) error {
	return nil
}
func (sm *TestSwarmManager) RemoveEndpoint(handle.Conn) error {
	return nil
}
func (sm *TestSwarmManager) Bisect() (SwarmManager, error) {
	return &TestSwarmManager{}, nil
}
func (sm *TestSwarmManager) Stitch(SwarmManager) error {
	return nil
}
func (sm *TestSwarmManager) Destroy() error {
	return nil
}

type FakeConn struct{}

func (fc *FakeConn) Read([]byte) (int, error)  { return 0, nil }
func (fc *FakeConn) Write([]byte) (int, error) { return 0, nil }
func (fc *FakeConn) Close() error              { return nil }
