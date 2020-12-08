package connector

import (
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"testing"
	"time"

	"github.com/arstevens/go-request/handle"
)

func TestConnector(t *testing.T) {
	totalRequests := 50
	requests := make([]TestConnectionRequest, totalRequests)
	conns := make([]FakeConn, totalRequests)

	for i := 0; i < totalRequests; i++ {
		logon := rand.Intn(100)%2 == 0
		originID := "/origin/" + strconv.Itoa(i%5)
		requestCode := rand.Intn(10)

		requests[i] = TestConnectionRequest{
			code:    requestCode,
			origin:  originID,
			logon:   logon,
			swarmID: "",
		}
		conns[i] = FakeConn{
			ip: net.IPv4(byte(rand.Intn(256)), byte(rand.Intn(256)),
				byte(rand.Intn(256)), byte(rand.Intn(256))),
		}
	}

	queueSize := 0
	verifier := TestIdentityVerifier{}
	connector := TestSwarmConnector{}

	handler := New(queueSize, &verifier, &connector)

	for i := 0; i < totalRequests; i++ {
		handler.AddJob(&requests[i], &conns[i])
		fmt.Println("")
	}
	time.Sleep(time.Second)
}

type TestIdentityVerifier struct{}

func (tv *TestIdentityVerifier) Analyze(ip net.IP, orig string, logon bool) bool {
	fmt.Printf("Verifying %s : %s : %t\n", ip.String(), orig, logon)
	return true
}

type TestSwarmConnector struct{}

func (tc *TestSwarmConnector) ProcessConnection(id string, code int, conn handle.Conn) error {
	fmt.Printf("Adding conn with code(%d) to swarm\n", code)
	return nil
}

type TestConnectionRequest struct {
	code    int
	origin  string
	logon   bool
	swarmID string
}

func (tr *TestConnectionRequest) GetRequestCode() int { return tr.code }
func (tr *TestConnectionRequest) GetOriginID() string { return tr.origin }
func (tr *TestConnectionRequest) IsLogOn() bool       { return tr.logon }
func (tr *TestConnectionRequest) GetSwarmID() string  { return tr.swarmID }

type FakeConn struct {
	ip net.IP
}

func (fc *FakeConn) GetIP() net.IP             { return fc.ip }
func (fc *FakeConn) Read([]byte) (int, error)  { return 0, nil }
func (fc *FakeConn) Write([]byte) (int, error) { return 0, nil }
func (fc *FakeConn) Close() error              { return nil }
