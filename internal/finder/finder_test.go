package finder

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
)

func TestFinder(t *testing.T) {
	totalSwarms := 10
	loads := make(map[string]int)
	prefLoads := make(map[string]int)
	swarms := make([]string, totalSwarms+1)
	for i := 0; i < totalSwarms; i++ {
		swarms[i] = "/dataspace/" + strconv.Itoa(i)
		loads[swarms[i]] = rand.Intn(1000) + 500
		prefLoads[swarms[i]] = rand.Intn(100)
	}
	prefLoads[swarms[totalSwarms-1]] = 0 //Division by zero case

	tracker := &testTracker{loads: loads, prefLoads: prefLoads}
	optimalFinder := New(tracker)
	for i := 0; i < totalSwarms; i++ {
		id := swarms[i]
		load := loads[id]
		prefLoad := prefLoads[id]
		optLoad := optimalFinder.GetBestSize(id)
		fmt.Printf("%s: Load->%d Pref_Load->%d Opt_Load->%d\n", id, load, prefLoad, optLoad)
	}
}

type testTracker struct {
	loads     map[string]int
	prefLoads map[string]int
}

func (tt *testTracker) GetLoad(id string) int {
	return tt.loads[id]
}
func (tt *testTracker) GetDebriefData(id string) interface{} {
	return tt.prefLoads[id]
}
