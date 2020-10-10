package localize

import (
	"net"

	"github.com/arstevens/go-request/handle"
)

// SwarmID is an alias for an integer used as an ID for a swarm handler
type SwarmID = int

/*ManagerMap maps SwarmIDs to request handlers which should be equipped
to handle swarm signaling requests*/
type SwarmHandlerMap interface {
	GetSwarmHandler(id SwarmID) (handle.RequestHandler, error)
}

/*SwarmMap takes in the content a device is requesting and its IP address
and attempts to locate the nearest with the requested data swarm for the
requester to connect to*/
type SwarmIDMap interface {
	GetSwarmID(dataID string, ipAddr net.IP) (SwarmID, error)
}
