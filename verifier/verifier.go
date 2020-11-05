package verifier

import "net"

type IdentityVerifier struct {
	registrationDB OriginDatabase
	connCache      ConnectionCache
}

func New(registrationDB OriginDatabase, connCache ConnectionCache) *IdentityVerifier {
	return &IdentityVerifier{
		registrationDB: registrationDB,
		connCache:      connCache,
	}
}

func (iv *IdentityVerifier) Analyze(ip net.IP, originID string, isLogOn bool) bool {
	valid := false
	if iv.registrationDB.IsRegistered(originID) {
		if isLogOn {
			valid = !iv.connCache.IsRecentlyConnected(ip)
			if valid {
				iv.connCache.NewConnection(ip)
			}
		} else {
			valid = true
		}
	}
	return valid
}
