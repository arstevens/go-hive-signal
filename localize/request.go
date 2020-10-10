package localize

import (
	"net"
)

/*DiscoverRequest describes a request object that
contains information about the requesters IP address
and the ID of the data they are requesting*/
type DiscoverRequest interface {
	GetIPAddress() net.IP
	GetDataID() string
}
