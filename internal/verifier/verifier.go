package verifier

import "net"

/*IdentityVerifier is an object that can check whether or
not a request to connect/disconnect from the network is valid*/
type IdentityVerifier struct {
	registrationDB OriginDatabase
	connCache      ConnectionCache
}

//New creates a new instance of IdentityVerifier
func New(registrationDB OriginDatabase, connCache ConnectionCache) *IdentityVerifier {
	return &IdentityVerifier{
		registrationDB: registrationDB,
		connCache:      connCache,
	}
}

//Analyze checks to see if a request is valid
func (iv *IdentityVerifier) Analyze(ip net.IP, originID string, isLogOn bool) bool {
	valid := false
	if iv.registrationDB.IsRegistered(originID) {
		if isLogOn && !iv.connCache.IsRecentlyConnected(ip) {
			iv.connCache.NewConnection(ip)
			valid = true
		} else if !iv.connCache.IsRecentlyDisconnected(ip) {
			iv.connCache.NewDisconnection(ip)
			valid = true
		}
	}
	return valid
}
