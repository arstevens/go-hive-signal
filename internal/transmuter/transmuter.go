package transmuter

import (
	"fmt"
	"log"

	"github.com/arstevens/go-request/handle"
)

const transmuterFailFormat = "Failed to process connection in SwarmTransmuter: %v"

/*SwarmTransmuter handles any commands that result in a change in
the makeup of a swarm*/
type SwarmTransmuter struct {
	swarmMap SwarmMap
	analyzer SwarmAnalyzer
}

//New creates a new SwarmTransmuter
func New(mapper SwarmMap, analyzer SwarmAnalyzer) *SwarmTransmuter {
	go pollForTransmutation(mapper, analyzer)
	return &SwarmTransmuter{
		swarmMap: mapper,
		analyzer: analyzer,
	}
}

//ProcessConnection processes a new request identified by 'code'
func (st *SwarmTransmuter) ProcessConnection(dataspaceID string, swarmConnect bool, conn handle.Conn) error {
	if swarmConnect {
		needyID, err := st.analyzer.GetMostNeedy()
		if err != nil {
			/*No swarms in need to a new endpoint so drop connection.
			In reality this only happens when the analyzer has yet to
			run it's first swarm analysis or if no swarms are registered*/
			log.Printf(transmuterFailFormat, err)
			return nil
		}
		m, err := st.swarmMap.GetSwarm(needyID)
		if err != nil {
			return fmt.Errorf(transmuterFailFormat, err)
		}
		manager := m.(SwarmManager)
		manager.AddEndpoint(conn)
	}
	return nil
}
