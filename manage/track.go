package manage

import "net"

/*MemberTracker describes an object that stores information
regarding the current state of a p2p endpoint*/
type MemberTracker interface {
	ModifyTrackingData(origin string, ip net.IP, active bool) error
}
