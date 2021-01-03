package transmuter

import (
	"fmt"
	"log"

	"github.com/arstevens/go-request/handle"
)

const (
	SwarmConnect = iota
	SwarmDisconnect
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
func (st *SwarmTransmuter) ProcessConnection(dataspaceID string, code int, conn handle.Conn) error {
	if code == SwarmConnect {
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
	} else if code == SwarmDisconnect {
		m, err := st.swarmMap.GetSwarm(dataspaceID)
		if err != nil {
			return fmt.Errorf(transmuterFailFormat, err)
		}
		manager := m.(SwarmManager)
		manager.RemoveEndpoint(conn)
	} else {
		return fmt.Errorf("Received invalid connection code in SwarmTransmuter")
	}
	return nil
}
