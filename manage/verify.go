package manage

import "net"

/*IdentityVerifier describes an object that can verify
whether or not an endpoint attempting to logon from ip
and claiming to come from origin is allowed to join a
swarm*/
type IdentityVerifier interface {
	Verify(origin string, ip net.IP) error
}
