package localize

import "net"

/*FrequencyManager describes an object that keeps track of
request frequency information to modify resource allocation*/
type FrequencyManager interface {
	IncrementFrequency(dataID string, ip net.IP)
}
