package transmuter

import (
	"fmt"

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
	sizeTracker SwarmSizeTracker
	swarmMap    SwarmMap
	analyzer    SwarmAnalyzer
}

//New creates a new SwarmTransmuter
func New(tracker SwarmSizeTracker, mapper SwarmMap, analyzer SwarmAnalyzer) *SwarmTransmuter {
	go pollForTransmutation(mapper, analyzer)
	return &SwarmTransmuter{
		sizeTracker: tracker,
		swarmMap:    mapper,
		analyzer:    analyzer,
	}
}

//ProcessConnection processes a new request identified by 'code'
func (st *SwarmTransmuter) ProcessConnection(swarmID string, code int, conn handle.Conn) error {
	if code == SwarmConnect {
		smallest, err := st.sizeTracker.GetSmallest()
		if err != nil {
			return fmt.Errorf(transmuterFailFormat, err)
		}
		manager, err := st.swarmMap.GetSwarmByID(smallest)
		if err != nil {
			return fmt.Errorf(transmuterFailFormat, err)
		}
		manager.AddEndpoint(conn)
	} else if code == SwarmDisconnect {
		manager, err := st.swarmMap.GetSwarmByID(swarmID)
		if err != nil {
			return fmt.Errorf(transmuterFailFormat, err)
		}
		manager.RemoveEndpoint(conn)
	} else {
		return fmt.Errorf("Received invalid connection code in SwarmTransmuter")
	}
	return nil
}
