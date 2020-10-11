package localize

import (
	"github.com/arstevens/go-request/handle"
)

// SwarmID is an alias for an integer used as an ID for a swarm handler
type SwarmID = int

/*SwarmMap takes in the content a device is requesting and
attempts to locate the nearest with the requested data swarm
for the requester to connect to*/
type SwarmMap interface {
	GetSwarmHandler(dataID string) (handle.RequestHandler, error)
}
