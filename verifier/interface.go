package verifier

import "net"

/*ConnectionCache describes an object that keeps track of
recently connected and disconnected device IP addresses*/
type ConnectionCache interface {
	IsRecentlyDisconnected(net.IP) bool
	IsRecentlyConnected(net.IP) bool
	NewConnection(net.IP)
	NewDisconnection(net.IP)
}

/*OriginDatabase describes an object that keeps track
of which points of origin are registered*/
type OriginDatabase interface {
	IsRegistered(string) bool
}
