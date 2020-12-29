package registrator

import (
	"fmt"
	"log"

	"github.com/arstevens/go-request/handle"
)

/*RegistrationHandler is an object that can handle requests that add a dataspace or
origin to the set of all supported by this signaling server*/
type RegistrationHandler struct {
	closed        bool
	requestStream chan<- handle.RequestPair
}

//New creates a new instance of RegistrationHandler
func New(size int, swarmMap SwarmMap, originReg OriginRegistrator) *RegistrationHandler {
	requestStream := make(chan handle.RequestPair, size)
	go processRequestStream(requestStream, swarmMap, originReg)
	return &RegistrationHandler{
		closed:        false,
		requestStream: requestStream,
	}
}

//AddJob adds a request to the job queue. AddJob hangs if the job queue is full
func (rh *RegistrationHandler) AddJob(request interface{}, conn handle.Conn) error {
	if rh.closed {
		return fmt.Errorf("Cannot add a job on a closed RegistrationHandler")
	}
	rh.requestStream <- handle.RequestPair{Request: request, Conn: conn}
	return nil
}

/*JobCapacity returns the max number of jobs a RegistrationHandler can hold
in its job queue*/
func (rh *RegistrationHandler) JobCapacity() int {
	return cap(rh.requestStream)
}

/*QueuedJobs returns the number of jobs currently queued in a RegistrationHandler
job queue */
func (rh *RegistrationHandler) QueuedJobs() int {
	return len(rh.requestStream)
}

/*Close closes a RegistrationHandler for work and returns an error if the
Handler is already closed*/
func (rh *RegistrationHandler) Close() error {
	if !rh.closed {
		rh.closed = true
		close(rh.requestStream)
		return nil
	}
	return fmt.Errorf("Cannot close a closed DataspaceHandler")
}

func processRequestStream(requestStream <-chan handle.RequestPair, swarmMap SwarmMap, originReg OriginRegistrator) {
	for {
		requestPair, ok := <-requestStream
		if !ok {
			return
		}
		request := requestPair.Request.(RegistrationRequest)

		err := handleRegistrationRequest(request, swarmMap, originReg)
		if err != nil {
			log.Println(err)
		}
	}
}

func handleRegistrationRequest(request RegistrationRequest, swarmMap SwarmMap, originReg OriginRegistrator) error {
	var err error
	if request.IsOrigin() {
		if request.IsAdd() {
			err = originReg.AddOrigin(request.GetDataField())
		} else {
			err = originReg.RemoveOrigin(request.GetDataField())
		}
	} else {
		if request.IsAdd() {
			err = swarmMap.AddSwarm(request.GetDataField())
		} else {
			err = swarmMap.RemoveSwarm(request.GetDataField())
		}
	}

	if err != nil {
		return fmt.Errorf("Failed to handle registration request in RegistrationHandler: %v", err)
	}
	return nil
}
