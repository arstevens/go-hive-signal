package analyzer

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func TestAnalyzer(t *testing.T) {
	rand.Seed(time.Now().Unix())

	totalDataspaces := 10
	sizeLimit := 100
	loadLimit := sizeLimit
	sizes := make(map[string]int)
	dupSizes := make(map[string]int)
	loads := make(map[string]int)

	fmt.Printf("Initial dataspace parameters\n----------------------------\n")
	for i := 0; i < totalDataspaces; i++ {
		dspace := "/dataspace/" + strconv.Itoa(i)
		sizes[dspace] = rand.Intn(sizeLimit)
		dupSizes[dspace] = sizes[dspace]
		loads[dspace] = rand.Intn(loadLimit)

		fmt.Printf("\t%s: (Size: %d) (Load: %d)\n", dspace, sizes[dspace], loads[dspace])
	}

	DistancePollTime = time.Millisecond
	tracker := &TestSwarmInfoTracker{sizes: sizes, loads: loads}
	finder := &TestOptimalSizeFinder{sizes: dupSizes}
	analyzer := New(tracker, finder)

	time.Sleep(DistancePollTime * 2)
	needyID, err := analyzer.GetMostNeedy()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("\nMost Needy\n----------\n\t%s\n", needyID)

	candidates, err := analyzer.CalculateCandidates()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("\nCandidates\n----------\n")
	for _, candidate := range candidates {
		fmt.Printf("\tTransferer: %s Transferee: %s Size: %d\n",
			candidate.GetTransfererID(), candidate.GetTransfereeID(), candidate.GetTransferSize())
	}

	fmt.Printf("\nFinal Distances\n---------------\n")
	for _, distance := range analyzer.matchDistances {
		fmt.Printf("\tDataspace: %s Distance: %d\n", distance.dataspace, distance.distance)
	}
}

type TestOptimalSizeFinder struct {
	sizes map[string]int
}

func (sf *TestOptimalSizeFinder) GetBestSize(id string) int {
	if size, ok := sf.sizes[id]; ok {
		return size
	}
	return -1
}

type TestSwarmInfoTracker struct {
	sizes map[string]int
	loads map[string]int
}

func (tt *TestSwarmInfoTracker) GetSize(id string) int {
	size, ok := tt.sizes[id]
	if !ok {
		return 0
	}
	return size
}
func (tt *TestSwarmInfoTracker) GetLoad(id string) int {
	load, ok := tt.loads[id]
	if !ok {
		return 0
	}
	return load
}
func (tt *TestSwarmInfoTracker) GetDataspaces() []string {
	dspaces := make([]string, 0, len(tt.sizes))
	for dspace, _ := range tt.sizes {
		dspaces = append(dspaces, dspace)
	}
	return dspaces
}
