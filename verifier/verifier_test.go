package verifier

import (
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"testing"
)

func TestVerifier(t *testing.T) {
	logOnProbability := 50
	totalOrigins := 5
	totalConnectors := 20
	totalRequests := 50

	origins := make([]string, totalOrigins)
	for i := 0; i < totalOrigins; i++ {
		origins[i] = "/origin/" + strconv.Itoa(i)
	}

	ipAddrs := make([]net.IP, totalConnectors)
	for i := 0; i < totalConnectors; i++ {
		ipAddrs[i] = net.IPv4(byte(rand.Intn(256)), byte(rand.Intn(256)),
			byte(rand.Intn(256)), byte(rand.Intn(256)))
	}

	connCache := TestConnectionCache{cache: make(map[string]bool),
		disCache: make(map[string]bool)}
	origDB := TestOriginDatabase{db: make(map[string]bool)}
	for i := 0; i < totalOrigins; i++ {
		origDB.db[origins[i]] = true
	}

	verifier := New(&origDB, &connCache)
	for i := 0; i < totalRequests; i++ {
		isLogOn := rand.Intn(100) < logOnProbability
		ip := ipAddrs[rand.Intn(totalConnectors)]
		origin := origins[rand.Intn(totalOrigins)]
		response := verifier.Analyze(ip, origin, isLogOn)
		fmt.Printf("Verification response: (%t) with parameters(%s, %s, %t)\n", response, ip.String(), origin, isLogOn)
	}
}

type TestConnectionCache struct {
	cache    map[string]bool
	disCache map[string]bool
}

func (tc *TestConnectionCache) IsRecentlyConnected(ip net.IP) bool {
	key := ip.String()
	_, ok := tc.cache[key]
	if ok {
		delete(tc.cache, key)
	}
	return ok
}
func (tc *TestConnectionCache) IsRecentlyDisconnected(ip net.IP) bool {
	key := ip.String()
	_, ok := tc.disCache[key]
	if ok {
		delete(tc.disCache, key)
	}
	return ok
}

func (tc *TestConnectionCache) NewConnection(ip net.IP) {
	key := ip.String()
	tc.cache[key] = true
}
func (tc *TestConnectionCache) NewDisconnection(ip net.IP) {
	tc.disCache[ip.String()] = true
}

type TestOriginDatabase struct {
	db map[string]bool
}

func (td *TestOriginDatabase) IsRegistered(id string) bool {
	_, ok := td.db[id]
	return ok
}
