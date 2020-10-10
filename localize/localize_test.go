package localize

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"testing"
	"time"

	"github.com/arstevens/go-request/handle"
)

func TestRequestLocalizer(t *testing.T) {
	swarmCount := 10
	swarmIDLimit := 100
	handlerCapacity := 10
	localizerQueueSize := 10
	totalJobs := 100

	ids := make([]SwarmID, swarmCount)
	handlers := make(map[SwarmID]handle.RequestHandler)
	for i := 0; i < swarmCount; i++ {
		randID := rand.Intn(swarmIDLimit)
		ids[i] = randID
		handlers[randID] = &TestHandler{capacity: handlerCapacity}
	}

	idMap := TestSwarmIDMap{ids: ids}
	handlerMap := TestSwarmHandlerMap{handlers: handlers}

	localizer, err := NewRequestLocalizer(localizerQueueSize, &idMap, &handlerMap)
	defer localizer.Close()

	for i := 0; i < totalJobs; i++ {
		job := TestLocalizeRequest{
			ip:     net.ParseIP("192.168.1.1"),
			dataID: "randID",
		}
		err = localizer.AddJob(&job)
		if err != nil {
			log.Println(err)
		}
	}
	time.Sleep(time.Second * 5)
}

type TestSwarmIDMap struct {
	ids []SwarmID
}

func (sm *TestSwarmIDMap) GetSwarmID(string, net.IP) (SwarmID, error) {
	return sm.ids[rand.Intn(len(sm.ids))], nil
}

type TestSwarmHandlerMap struct {
	handlers map[SwarmID]handle.RequestHandler
}

func (sh *TestSwarmHandlerMap) GetSwarmHandler(id SwarmID) (handle.RequestHandler, error) {
	handler, ok := sh.handlers[id]
	if !ok {
		return nil, fmt.Errorf("Invalid swarm id")
	}
	return handler, nil
}

type TestLocalizeRequest struct {
	ip     net.IP
	dataID string
}

func (lr *TestLocalizeRequest) GetIPAddress() net.IP { return lr.ip }
func (lr *TestLocalizeRequest) GetDataID() string    { return lr.dataID }

type TestHandler struct {
	capacity int
	queued   int
}

func (h *TestHandler) AddJob(interface{}) error {
	if h.queued != h.capacity {
		h.queued++
	}
	fmt.Printf("Received new job! %d/%d spots used\n", h.queued, h.capacity)
	return nil
}

func (h *TestHandler) JobCapacity() int {
	return h.capacity
}

func (h *TestHandler) QueuedJobs() int {
	return h.queued
}

func (h *TestHandler) Close() error {
	return nil
}
