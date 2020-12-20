package integration_tests

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/arstevens/go-hive-signal/internal/gateway"
	"github.com/arstevens/go-hive-signal/internal/manager"
	"github.com/arstevens/go-hive-signal/internal/negotiator"
	"github.com/arstevens/go-hive-signal/internal/tracker"
)

func TestIntegrationOne(t *testing.T) {
	gatewayActiveSize, gatewayInactiveSize := 10, 20
	sTracker := tracker.New()

	totalSwarms := 10
	fmt.Printf("Stage One: Creating %d swarms...\n", totalSwarms)
	swarms := make([]*manager.SwarmManager, totalSwarms)
	gateways := make([]*gateway.SwarmGateway, totalSwarms)

	//fakeNegotiate := func(manager.Conn, manager.Conn) error { return nil }
	negotiator.UnmarshalMessage = unmarshalMessage
	for i := 0; i < totalSwarms; i++ {
		sGateway := gateway.New(gatewayActiveSize, gatewayInactiveSize)
		swarms[i] = manager.New("/swarm/"+strconv.Itoa(i), sGateway, negotiator.RoundtripLimitedNegotiate, sTracker)
		gateways[i] = sGateway
	}

	endpointCache := make(map[int][]string)
	endpointsPerSwarm := gatewayActiveSize + gatewayInactiveSize
	manager.ChangeTriggerLimit = 5
	fmt.Printf("Stage Two: Adding %d endpoints per swarm...\n", endpointsPerSwarm)
	for i := 0; i < totalSwarms; i++ {
		endpointCache[i] = make([]string, 0)
		for j := 0; j < endpointsPerSwarm; j++ {
			addr := "/endpoint/" + strconv.Itoa(i*endpointsPerSwarm+j)
			conn := newFakeConn(addr)
			err := swarms[i].AddEndpoint(conn)
			if err != nil {
				t.Fatal(err)
			}
			endpointCache[i] = append(endpointCache[i], addr)
		}
	}

	totalPairs := 800
	totalTransmutations := 800
	fmt.Printf("Stage Three: Running %d pairings and %d transmutations in simultaneous...\n", totalPairs, totalTransmutations)
	pairingsDone := make(chan struct{})
	transmutationsDone := make(chan struct{})

	gateway.DialEndpoint = func(addr string) (gateway.Conn, error) {
		return newFakeConn(addr), nil
	}

	swarmArrayMutex := &sync.Mutex{}
	pairingsFunc := func() {
		defer close(pairingsDone)
		for i := 0; i < totalPairs; i++ {
			conn := newFakeConn("/requester/" + strconv.Itoa(i))
			swarmArrayMutex.Lock()
			swarm := swarms[rand.Intn(len(swarms))]
			swarmArrayMutex.Unlock()
			err := swarm.AttemptToPair(conn)
			if err != nil && err.Error()[0] != 'F' { //Error starting in F is inactive queue empty error
				t.Fatal(err)
			}
		}
	}
	go pairingsFunc()

	go func() {
		defer close(transmutationsDone)
		for i := 0; i < totalTransmutations; i++ {
			mutType := rand.Intn(2)
			if mutType == 0 { //AddEndpoint
				addr := "/endpoint/NEW_" + strconv.Itoa(i)
				conn := newFakeConn(addr)
				swarmIdx := rand.Intn(len(swarms))
				swarm := swarms[swarmIdx]
				err := swarm.AddEndpoint(conn)
				if err != nil {
					t.Fatal(err)
				}
				endpointCache[swarmIdx] = append(endpointCache[swarmIdx], addr)
			} else if mutType == 1 { //RemoveEndpoint
				swarmIdx := rand.Intn(len(swarms))
				endpointIdx := rand.Intn(len(endpointCache[swarmIdx]))

				addr := endpointCache[swarmIdx][endpointIdx]
				endpointCache[swarmIdx] = append(endpointCache[swarmIdx][:endpointIdx], endpointCache[swarmIdx][endpointIdx+1:]...)
				conn := newFakeConn(addr)

				swarm := swarms[swarmIdx]
				err := swarm.RemoveEndpoint(conn)
				if err != nil {
					t.Fatal(err)
				}
			}
		}
	}()

	<-pairingsDone
	<-transmutationsDone

	fmt.Printf("Stage Four: Checking validity of tracked swarm sizes...\n")
	for i := 0; i < totalSwarms; i++ {
		id := "/swarm/" + strconv.Itoa(i)
		trackedSize := sTracker.GetSize(id)
		realSize := gateways[i].GetTotalEndpoints()
		difference := int(math.Abs(float64(trackedSize - realSize)))
		if difference > manager.ChangeTriggerLimit {
			t.Fatal(fmt.Errorf("Size difference between real size and tracked size is greater than trigger limit"))
		}
	}

	totalBisectMerges := 40
	fmt.Printf("Stage Five: Running %d pairings and %d bisections/merges in simultaneous...\n", totalPairs, totalBisectMerges)

	pairingsDone = make(chan struct{})
	transmutationsDone = make(chan struct{})
	go pairingsFunc()
	go func() {
		defer close(transmutationsDone)
		for i := 0; i < totalBisectMerges; i++ {
			swarmArrayMutex.Lock()
			canStitch := len(swarms) > 1
			swarmIdx := rand.Intn(len(swarms))
			swarmIdx2 := (swarmIdx + 1) % len(swarms)
			swarm1 := swarms[swarmIdx]
			var swarm2 *manager.SwarmManager
			if canStitch {
				swarm2 = swarms[swarmIdx2]
			}
			swarmArrayMutex.Unlock()

			mutType := rand.Intn(2)
			if mutType == 0 { //Bisect
				m, err := swarm1.Bisect()
				if err != nil {
					t.Fatal(err)
				}
				manager := m.(*manager.SwarmManager)
				manager.SetID("/swarms/" + strconv.Itoa(totalSwarms+i))
				swarmArrayMutex.Lock()
				swarms = append(swarms, manager)
				swarmArrayMutex.Unlock()
			} else if mutType == 1 && canStitch { //Stitch
				err := swarm1.Stitch(swarm2)
				if err != nil {
					t.Fatal(err)
				}
				swarmArrayMutex.Lock()
				swarms = append(swarms[:swarmIdx2], swarms[swarmIdx2+1:]...)
				swarmArrayMutex.Unlock()
			}
		}
	}()

	<-pairingsDone
	<-transmutationsDone

	fmt.Printf("Stage Six: Outputing tracker swarm sizes...\n")
	for i := 0; i < len(swarms); i++ {
		id := swarms[i].GetID()
		trackedSize := sTracker.GetSize(id)
		fmt.Printf("\t%s: %d\n", id, trackedSize)
	}
}

type FakeConn struct {
	addr   string
	buf    []byte
	closed bool
}

func newFakeConn(addr string) *FakeConn {
	return &FakeConn{
		addr:   addr,
		buf:    make([]byte, 0),
		closed: false,
	}
}

func (fc *FakeConn) Read(b []byte) (int, error) {
	return len(b), nil
}

func (fc *FakeConn) Write(b []byte) (int, error) {
	return len(b), nil
}

type fakeNegotiateMessage struct{}

func (nm *fakeNegotiateMessage) IsAccepted() bool {
	return true
}
func unmarshalMessage(b []byte) (interface{}, error) {
	return &fakeNegotiateMessage{}, nil
}

/*

func (fc *FakeConn) Read(b []byte) (int, error) {
	i := 0
	for ; i < len(b) && i < len(fc.buf); i++ {
		b[i] = fc.buf[i]
	}
	fc.buf = fc.buf[i:]
	return i, nil
}

func (fc *FakeConn) Write(b []byte) (int, error) {
	for i := 0; i < len(b); i++ {
		fc.buf = append(fc.buf, b[i])
	}
	return len(b), nil
}
*/

func (fc *FakeConn) Close() error       { fc.closed = true; return nil }
func (fc *FakeConn) IsClosed() bool     { return fc.closed }
func (fc *FakeConn) GetAddress() string { return fc.addr }
