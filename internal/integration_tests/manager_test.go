package integration_tests

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/arstevens/go-hive-signal/internal/analyzer"
	"github.com/arstevens/go-hive-signal/internal/gateway"
	"github.com/arstevens/go-hive-signal/internal/localizer"
	"github.com/arstevens/go-hive-signal/internal/manager"
	"github.com/arstevens/go-hive-signal/internal/mapper"
	"github.com/arstevens/go-hive-signal/internal/tracker"
	"github.com/arstevens/go-hive-signal/internal/transmuter"
)

func TestRequestLocalizerSubtree(t *testing.T) {
	fmt.Printf("\nREQUEST LOCALIZER SUBTREE\n----------------------\n")
	//Prepare environment
	analyzer.OptimalSizeForLoad = func(s int) int {
		if s%2 == 0 {
			return s * 20
		}
		return s
	}
	analyzer.DistancePollTime = time.Millisecond * 10
	tracker.FrequencyCalculationPeriod = time.Second //time.Millisecond * 50
	transmuter.PollPeriod = time.Second
	gateway.DialEndpoint = func(addr string) (manager.Conn, error) {
		return &FakeConn{
			addr:   addr,
			closed: false,
		}, nil
	}
	negotiate := func(a manager.Conn, b manager.Conn) error { return nil }

	logName := "test.log"
	var err error
	logOut, err := os.Create(logName)
	if err != nil {
		t.Fatal(fmt.Errorf("Failed to open log file %s. Exiting...\n", logName))
	}
	defer logOut.Close()
	log.SetOutput(logOut)

	//Create all needed instances
	fmt.Printf("Creating simulation instances...\n")
	historyLength := 10
	infoTracker := tracker.New(historyLength)

	activeSize := 10
	inactiveSize := 20
	gatewayGen := gateway.NewGenerator(activeSize, inactiveSize)
	managerGen := manager.NewGenerator(gatewayGen, negotiate, infoTracker)

	swarmMap := mapper.New(managerGen)
	dataRequestAnalyzer := analyzer.New(infoTracker)
	swarmTransmuter := transmuter.New(swarmMap, dataRequestAnalyzer)

	localizerSize := 10
	requestLocalizer := localizer.New(localizerSize, swarmMap, infoTracker)

	//Register dataspaces
	totalDataspaces := 10
	fmt.Printf("Registering %d dataspaces...\n", totalDataspaces)
	dataspaces := make([]string, totalDataspaces)
	endpointSet := make(map[string]map[string]bool)
	for i := 0; i < totalDataspaces; i++ {
		dataspaces[i] = fmt.Sprintf("/dataspace/%d", i)
		endpointSet[dataspaces[i]] = make(map[string]bool)
		err = swarmMap.AddSwarm(dataspaces[i])
		if err != nil {
			t.Fatal(err)
		}
	}

	totalRuntime := time.Second * 5

	//Start endpoint additions/removals
	rand.Seed(time.Now().UnixNano())
	transDone := make(chan struct{})
	transmutationFrequency := time.Millisecond
	transmutationIterations := int(totalRuntime / transmutationFrequency)

	fmt.Printf("Starting endpoint additions...\n")
	newEndpointFrequency := 0.75
	go func() {
		defer close(transDone)
		var err error
		outOf := 100
		newEndpointFrequencyLimit := int(float64(outOf) * newEndpointFrequency)
		for i := 0; i < transmutationIterations; i++ {
			time.Sleep(transmutationFrequency)

			if rand.Intn(outOf) < newEndpointFrequencyLimit { //Add new endpoint
				dataspace := dataspaces[rand.Intn(totalDataspaces)]
				conn := &FakeConn{
					addr: fmt.Sprintf("%d.%d.%d.%d", rand.Intn(256), rand.Intn(256),
						rand.Intn(256), rand.Intn(256)),
					closed: false,
				}
				endpointSet[dataspace][conn.addr] = true
				err = swarmTransmuter.ProcessConnection(dataspace, transmuter.SwarmConnect, conn)
				if err != nil {
					fmt.Printf("%v\n", err)
					t.Fatal(err)
				}
			} else { //remove old endpoint
				dataspace := dataspaces[rand.Intn(totalDataspaces)]
				endpoint := getEndpointAndRemove(dataspace, endpointSet)
				if endpoint == "" {
					continue
				}

				conn := &FakeConn{
					addr:   endpoint,
					closed: false,
				}
				err = swarmTransmuter.ProcessConnection(dataspace, transmuter.SwarmDisconnect, conn)
				if err != nil {
					t.Fatal(err)
				}
			}

		}
	}()

	//Start data requests
	requestDone := make(chan struct{})
	requestFrequency := time.Millisecond
	requestIterations := int((totalRuntime - tracker.FrequencyCalculationPeriod) / requestFrequency)
	fmt.Printf("Starting pairing requests...\n")
	go func() {
		defer close(requestDone)
		time.Sleep(tracker.FrequencyCalculationPeriod)
		var err error
		for i := 0; i < requestIterations; i++ {
			time.Sleep(requestFrequency)

			conn := &FakeConn{
				addr: fmt.Sprintf("%d.%d.%d.%d", rand.Intn(256), rand.Intn(256),
					rand.Intn(256), rand.Intn(256)),
				closed: false,
			}
			job := TestLocalizeRequest(dataspaces[rand.Intn(totalDataspaces)])
			err = requestLocalizer.AddJob(&job, conn)
			if err != nil {
				t.Fatal(err)
			}
		}
	}()
	<-requestDone
	<-transDone

	//State output
	fmt.Printf("Outputing server state...\n")
	fmt.Printf("Tracker output\n")
	fmt.Printf("%v\n", infoTracker.GetDataspaces())
	for _, dspace := range dataspaces {
		load := infoTracker.GetLoad(dspace)
		tsize := infoTracker.GetSize(dspace)
		m, _ := swarmMap.GetSwarm(dspace)
		man := m.(*manager.SwarmManager)
		size := len(man.GetEndpoints())
		fmt.Printf("\t%s Load->%d TSize->%d Size->%d\n", dspace, load, tsize, size)
	}
}

func getEndpointAndRemove(dspace string, m map[string]map[string]bool) string {
	endpoints := m[dspace]
	length := len(endpoints)
	if length == 0 {
		return ""
	}
	idx := rand.Intn(length)
	endpoint := ""
	i := 0
	for e, _ := range endpoints {
		if i == idx {
			endpoint = e
			break
		}
		i++
	}

	if endpoint == "" {
		return endpoint
	}
	delete(endpoints, endpoint)
	return endpoint
}

type TestLocalizeRequest string

func (lr *TestLocalizeRequest) GetDataspace() string { return string(*lr) }

type FakeConn struct {
	addr   string
	closed bool
}

func (fc *FakeConn) Read([]byte) (int, error)  { return 0, nil }
func (fc *FakeConn) Write([]byte) (int, error) { return 0, nil }
func (fc *FakeConn) Close() error              { fc.closed = true; return nil }
func (fc *FakeConn) IsClosed() bool            { return fc.closed }
func (fc *FakeConn) GetAddress() string        { return fc.addr }
