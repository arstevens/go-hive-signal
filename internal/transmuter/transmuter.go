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
	gateway     SwarmGateway
}

//New creates a new SwarmTransmuter
func New(tracker SwarmSizeTracker, mapper SwarmMap, analyzer SwarmAnalyzer,
	gateway SwarmGateway) *SwarmTransmuter {
	go pollForTransmutation(mapper, gateway, analyzer)
	return &SwarmTransmuter{
		sizeTracker: tracker,
		swarmMap:    mapper,
		analyzer:    analyzer,
		gateway:     gateway,
	}
}

//ProcessConnection processes a new request identified by 'code'
func (st *SwarmTransmuter) ProcessConnection(swarmID string, code int, conn handle.Conn) error {
	if code == SwarmConnect {
		smallest, err := st.sizeTracker.GetSmallest()
		if err != nil {
			return fmt.Errorf(transmuterFailFormat, err)
		}
		err = st.gateway.AddEndpoint(smallest, conn)
		if err != nil {
			return fmt.Errorf(transmuterFailFormat, err)
		}
		st.sizeTracker.Increment(smallest)
	} else if code == SwarmDisconnect {
		st.sizeTracker.Decrement(swarmID)
	} else {
		return fmt.Errorf("Received invalid connection code in SwarmTransmuter")
	}
	return nil
}
