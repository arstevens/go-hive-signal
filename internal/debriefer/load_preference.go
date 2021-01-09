package debriefer

import (
	"encoding/binary"
	"io"
	"log"

	"github.com/arstevens/go-hive-signal/internal/tracker"
)

/*LoadPreferenceDebrief reads a single 32-bit integer value from
a connection that represents the connections preferred load*/
func LoadPreferenceDebrief(conn io.Reader) interface{} {
	var debriefValue int32
	err := binary.Read(conn, binary.BigEndian, &debriefValue)
	if err != nil {
		log.Printf("Failed to debrief in gateway.debriefConnection(): %v", err)
		return -1
	}
	return int(debriefValue)
}

//LPSEGenerator implements tracker.SwarmEngineGenerator
type LPSEGenerator struct {
	historyLength int
}

//NewLPSEGenerator returns a new instance of LPSEGenerator
func NewLPSEGenerator(historyLength int) *LPSEGenerator {
	return &LPSEGenerator{historyLength}
}

//New creates a new instance of LoadPreferenceStorageEngine
func (g *LPSEGenerator) New() tracker.StorageEngine {
	return &LoadPreferenceStorageEngine{
		loadTracker: tracker.NewLoadTracker(g.historyLength),
	}
}

/*LoadPreferenceStorageEngine implements tracker.StorageEngine
and keeps track of a load preferrence history and returns
the calculated average load preferrence*/
type LoadPreferenceStorageEngine struct {
	loadTracker *tracker.SwarmLoadTracker
}

//AddDatapoint adds the 'debrief' load preferrence to the storage engine
func (lp *LoadPreferenceStorageEngine) AddDatapoint(debrief interface{}) {
	lp.loadTracker.AddFrequencyDatapoint(debrief.(int))
}

//GetData retrieves the average load preferrence
func (lp *LoadPreferenceStorageEngine) GetData() interface{} {
	return lp.loadTracker.CalculateAverageFrequency()
}
