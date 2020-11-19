package cache

import (
	"net"
	"sync"
	"time"
)

var (
	//ConnectionTTL is the time-to-live of a connection record
	ConnectionTTL = time.Minute * 5
	//DisconnectionTTL is the time-to-live of a disconnection record
	DisconnectionTTL = time.Minute
)

/*GarbageCollectionPeriod is the frequency at which the cache is
purged of old records*/
var GarbageCollectionPeriod = time.Minute

/*ConnectionCache is an object that keeps track of connect/disconnect
requests*/
type ConnectionCache struct {
	mutex           *sync.Mutex
	connectCache    map[string]time.Time
	disconnectCache map[string]time.Time
}

//New creates a new ConnectionCache
func New() *ConnectionCache {
	cache := ConnectionCache{
		mutex:           &sync.Mutex{},
		connectCache:    make(map[string]time.Time),
		disconnectCache: make(map[string]time.Time),
	}
	go pollForTimedOutRecords(cache.mutex, []map[string]time.Time{
		cache.connectCache,
		cache.disconnectCache,
	})
	return &cache
}

//IsRecentlyDisconnected returns whether or not 'ip' recently disconnected
func (cc *ConnectionCache) IsRecentlyDisconnected(ip net.IP) bool {
	cc.mutex.Lock()
	defer cc.mutex.Unlock()

	ipKey := ip.String()
	ttl, ok := cc.disconnectCache[ipKey]
	if !ok {
		return false
	}
	if ttl.Before(time.Now()) {
		delete(cc.disconnectCache, ipKey)
		return false
	}
	return true
}

//IsRecentlyConnected returns whether or not 'ip' recently connected
func (cc *ConnectionCache) IsRecentlyConnected(ip net.IP) bool {
	cc.mutex.Lock()
	defer cc.mutex.Unlock()

	ipKey := ip.String()
	ttl, ok := cc.connectCache[ipKey]
	if !ok {
		return false
	}
	if ttl.Before(time.Now()) {
		delete(cc.connectCache, ipKey)
		return false
	}
	return true
}

//NewDisconnection adds a new disconnection record for 'ip'
func (cc *ConnectionCache) NewDisconnection(ip net.IP) {
	cc.mutex.Lock()
	defer cc.mutex.Unlock()

	ttl := time.Now().Add(DisconnectionTTL)
	cc.disconnectCache[ip.String()] = ttl
}

//NewConnection adds a new connection record for 'ip'
func (cc *ConnectionCache) NewConnection(ip net.IP) {
	cc.mutex.Lock()
	defer cc.mutex.Unlock()

	ttl := time.Now().Add(ConnectionTTL)
	cc.connectCache[ip.String()] = ttl
}
