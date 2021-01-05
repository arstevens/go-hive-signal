package integration_tests

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"testing"
	"time"

	"github.com/arstevens/go-hive-signal/internal/analyzer"
	"github.com/arstevens/go-hive-signal/internal/cache"
	"github.com/arstevens/go-hive-signal/internal/connector"
	"github.com/arstevens/go-hive-signal/internal/gateway"
	"github.com/arstevens/go-hive-signal/internal/localizer"
	"github.com/arstevens/go-hive-signal/internal/manager"
	"github.com/arstevens/go-hive-signal/internal/mapper"
	"github.com/arstevens/go-hive-signal/internal/register"
	"github.com/arstevens/go-hive-signal/internal/registrator"
	"github.com/arstevens/go-hive-signal/internal/tracker"
	"github.com/arstevens/go-hive-signal/internal/transmuter"
	"github.com/arstevens/go-hive-signal/internal/verifier"
)

/*
func TestIdentityVerifier(t *testing.T) {
	endpointRegister, err := register.New()
	if err != nil {
		t.Fatal(err)
	}

	totalOrigins := 10
	origins := make([]string, totalOrigins)
	for i := 0; i < totalOrigins; i++ {
		id := "/origin/" + strconv.Itoa(i)
		err = endpointRegister.AddOrigin(id)
		if err != nil {
			t.Fatal(err)
		}
		origins[i] = id
	}

	cache.GarbageCollectionPeriod = time.Millisecond * 3
	cache.ConnectionTTL = time.Millisecond * 4
	cache.DisconnectionTTL = time.Millisecond * 4
	connectionCache := cache.New()
	identityVerifier := verifier.New(endpointRegister, connectionCache)

	totalIPs := 10
	ips := make([]net.IP, totalIPs)
	for i := 0; i < totalIPs; i++ {
		ip := net.IPv4(byte(rand.Intn(256)), byte(rand.Intn(256)),
			byte(rand.Intn(256)), byte(rand.Intn(256)))
		ips[i] = ip
	}

	totalConnections := 20
	for i := 0; i < totalConnections; i++ {
		time.Sleep(time.Millisecond)
		idx := rand.Intn(totalIPs)
		ip := ips[idx]

		isLogOn := rand.Intn(10) > 5
		isValid := identityVerifier.Analyze(ip, origins[0], isLogOn)
		fmt.Printf("IsLogOn: %t, IP: %s, IsValid: %t\n", isLogOn, ip.String(), isValid)
	}

	for _, origin := range origins {
		err = endpointRegister.RemoveOrigin(origin)
		if err != nil {
			log.Printf("%v\n", err)
		}
	}
}
*/

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

	requestBufferSize := 10
	requestLocalizer := localizer.New(requestBufferSize, swarmMap, infoTracker)

	endpointRegister, err := register.New()
	if err != nil {
		t.Fatal(err)
	}
	registrationHandler := registrator.New(requestBufferSize, swarmMap, endpointRegister)

	cache.GarbageCollectionPeriod = time.Millisecond * 50
	cache.ConnectionTTL = time.Second
	cache.DisconnectionTTL = time.Second
	connectionCache := cache.New()
	identityVerifier := verifier.New(endpointRegister, connectionCache)
	connectionHandler := connector.New(requestBufferSize, identityVerifier, swarmTransmuter)

	//Register dataspaces
	totalDataspaces := 10
	fmt.Printf("Registering %d dataspaces...\n", totalDataspaces)
	dataspaces := make([]string, totalDataspaces)
	endpointSet := make(map[string]map[string]bool)
	for i := 0; i < totalDataspaces; i++ {
		dataspaces[i] = fmt.Sprintf("/dataspace/%d", i)
		endpointSet[dataspaces[i]] = make(map[string]bool)

		registrationRequest := &TestRegistrationRequest{
			add:    true,
			origin: false,
			field:  dataspaces[i],
		}

		err = registrationHandler.AddJob(registrationRequest, &FakeConn{})
		if err != nil {
			t.Fatal(err)
		}
	}

	totalOrigins := 10
	fmt.Printf("Registering %d points of origin...\n", totalOrigins)
	origins := make([]string, totalOrigins)
	for i := 0; i < totalOrigins; i++ {
		origins[i] = fmt.Sprintf("/origin/%d", i)

		registrationRequest := &TestRegistrationRequest{
			add:    true,
			origin: true,
			field:  origins[i],
		}
		err = registrationHandler.AddJob(registrationRequest, &FakeConn{})
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

				connectionRequest := &TestConnectionRequest{
					code:     transmuter.SwarmConnect,
					swarmID:  dataspace,
					originID: origins[0],
					isLogOn:  true,
				}
				err = connectionHandler.AddJob(connectionRequest, conn)
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
				disconnectionRequest := &TestConnectionRequest{
					code:     transmuter.SwarmDisconnect,
					swarmID:  dataspace,
					originID: origins[0],
					isLogOn:  false,
				}
				err = connectionHandler.AddJob(disconnectionRequest, conn)
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
	for _, dspace := range dataspaces {
		load := infoTracker.GetLoad(dspace)
		tsize := infoTracker.GetSize(dspace)
		m, _ := swarmMap.GetSwarm(dspace)
		man := m.(*manager.SwarmManager)
		size := len(man.GetEndpoints())
		fmt.Printf("\t%s Load->%d TSize->%d Size->%d\n", dspace, load, tsize, size)
	}

	var fatalErr error
	for _, origin := range origins {
		registrationRequest := &TestRegistrationRequest{
			add:    false,
			origin: true,
			field:  origin,
		}
		err = registrationHandler.AddJob(registrationRequest, &FakeConn{})
		if err != nil {
			fatalErr = err
			log.Printf("%v\n", err)
		}
	}
	if fatalErr != nil {
		t.Fatal(err)
	}
	for _, dspace := range dataspaces {
		registrationRequest := &TestRegistrationRequest{
			add:    false,
			origin: false,
			field:  dspace,
		}
		err = registrationHandler.AddJob(registrationRequest, &FakeConn{})
		if err != nil {
			t.Fatal(err)
		}
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

type TestConnectionRequest struct {
	code     int
	swarmID  string
	originID string
	isLogOn  bool
}

func (cr *TestConnectionRequest) GetRequestCode() int { return cr.code }
func (cr *TestConnectionRequest) GetSwarmID() string  { return cr.swarmID }
func (cr *TestConnectionRequest) GetOriginID() string { return cr.originID }
func (cr *TestConnectionRequest) IsLogOn() bool       { return cr.isLogOn }

type TestRegistrationRequest struct {
	add    bool
	origin bool
	field  string
}

func (rr *TestRegistrationRequest) IsAdd() bool          { return rr.add }
func (rr *TestRegistrationRequest) IsOrigin() bool       { return rr.origin }
func (rr *TestRegistrationRequest) GetDataField() string { return rr.field }

type TestOriginRegistrator struct{}

func (or *TestOriginRegistrator) AddOrigin(string) error    { return nil }
func (or *TestOriginRegistrator) RemoveOrigin(string) error { return nil }

type TestLocalizeRequest string

func (lr *TestLocalizeRequest) GetDataspace() string { return string(*lr) }

type FakeConn struct {
	addr   string
	closed bool
}

func (fc *FakeConn) GetIP() net.IP             { return net.ParseIP(fc.addr) }
func (fc *FakeConn) Read([]byte) (int, error)  { return 0, nil }
func (fc *FakeConn) Write([]byte) (int, error) { return 0, nil }
func (fc *FakeConn) Close() error              { fc.closed = true; return nil }
func (fc *FakeConn) IsClosed() bool            { return fc.closed }
func (fc *FakeConn) GetAddress() string        { return fc.addr }
