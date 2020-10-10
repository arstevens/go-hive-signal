package manage

import "net"

/*ConnectionRequest describes a request object that contains
information regarding the requesters application of origin
and its IP address. This information is used to verify the
validity of the request and to log the requester onto the
network*/
type ConnectionRequest interface {
	GetOrigin() string
	GetIPAddress() net.IP
	GetActivityStatus() bool
}
