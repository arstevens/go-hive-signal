package manage

import "net"

/*MemberTracker describes an object that stores information
regarding the current state of a p2p endpoint*/
type MemberTracker interface {
	StopTracking(ip net.IP) error
	ModifyTrackingData(origin string, ip net.IP, active bool) error
}
