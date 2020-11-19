package cache

import (
	"fmt"
	"math/rand"
	"net"
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	fmt.Printf("---------------CONNECTION CACHE TEST------------------\n")
	ConnectionTTL = time.Second
	DisconnectionTTL = time.Second
	GarbageCollectionPeriod = time.Second

	totalIPs := 10
	totalPings := 30
	ips := make([]net.IP, totalIPs)
	for i := 0; i < totalIPs; i++ {
		ips[i] = net.IPv4(byte(rand.Intn(256)), byte(rand.Intn(256)),
			byte(rand.Intn(256)), byte(rand.Intn(256)))
	}

	cache := New()
	iterations := 3
	for j := 0; j < iterations; j++ {
		fmt.Printf("------------------\nIteration %d\n------------------\n", j)
		for i := 0; i < totalPings; i++ {
			ip := ips[rand.Intn(totalIPs)]
			randNum := rand.Intn(100)
			if randNum < 25 {
				cache.NewConnection(ip)
				fmt.Printf("(%s LOGGED CONNECTION)\n", ip.String())
			} else if randNum < 50 {
				fmt.Printf("(%s LOGGED DISCONNECTION)\n", ip.String())
				cache.NewDisconnection(ip)
			} else if randNum < 75 {
				fmt.Printf("(%s)[CONNECTION STATUS] = %t\n", ip.String(), cache.IsRecentlyConnected(ip))
			} else {
				fmt.Printf("(%s)[DISCONNECTION STATUS] = %t\n", ip.String(), cache.IsRecentlyDisconnected(ip))
			}
		}
		fmt.Printf("[CACHE STATUS PRE-COLLECTION]")
		fmt.Printf("Connect Cache -> ")
		printKeys(cache.connectCache)
		fmt.Printf("Disconnect Cache -> ")
		printKeys(cache.disconnectCache)
		time.Sleep(GarbageCollectionPeriod * 2)
		fmt.Printf("[CACHE STATUS POST-COLLECTION]")
		fmt.Printf("Connect Cache -> ")
		printKeys(cache.connectCache)
		fmt.Printf("Disconnect Cache -> ")
		printKeys(cache.disconnectCache)
	}
}

func printKeys(m map[string]time.Time) {
	fmt.Printf("{\n")
	for key, _ := range m {
		fmt.Printf("\t%s,\n", key)
	}
	fmt.Printf("}\n")
}
