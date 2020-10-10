package manage

import (
	"fmt"
	"math/rand"
	"net"
	"testing"
	"time"
)

func TestMemberManager(t *testing.T) {
	queueSize := 10
	totalJobs := 20
	totalDataspaces := 10
	jobs := make([]TestRequest, totalJobs)
	for i := 0; i < totalJobs; i++ {
		jobs[i] = TestRequest{
			origin: fmt.Sprintf("/dataspace/%d", rand.Intn(totalDataspaces)),
			ip:     net.ParseIP(fmt.Sprintf("192.168.1.%d", rand.Intn(256))),
			status: rand.Intn(100) < 50,
		}
	}

	verifier := TestVerifier{}
	tracker := TestTracker{}
	handler, _ := NewMemberManager(queueSize, &verifier, &tracker)
	defer handler.Close()

	for i := 0; i < totalJobs; i++ {
		handler.AddJob(&jobs[i])
	}
	time.Sleep(time.Second * 5)
}

type TestVerifier struct{}

func (v *TestVerifier) Verify(s string, i net.IP) error {
	fmt.Printf("Verifying (%s, %v)\n", s, i)
	return nil
}

type TestTracker struct{}

func (t *TestTracker) ModifyTrackingData(s string, i net.IP, b bool) error {
	fmt.Printf("Modifying tracking data for (%s, %v) to status: %t\n", s, i, b)
	return nil
}

type TestRequest struct {
	origin string
	ip     net.IP
	status bool
}

func (r *TestRequest) GetOrigin() string       { return r.origin }
func (r *TestRequest) GetIPAddress() net.IP    { return r.ip }
func (r *TestRequest) GetActivityStatus() bool { return r.status }
