package verifier

import "net"

type ConnectionCache interface {
	IsRecentlyConnected(net.IP) bool
	NewConnection(net.IP)
}

type OriginDatabase interface {
	IsRegistered(string) bool
}
