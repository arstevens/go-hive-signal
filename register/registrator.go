package register

import (
	"fmt"
	"log"

	"github.com/arstevens/go-request/handle"
)

/*DataspaceHandler is an object that can handle requests that add a dataspace to
the set of all supported dataspaces for a signaling server*/
type RegistrationHandler struct {
	closed        bool
	requestStream chan<- handle.RequestPair
}

func New(size int, swarmMap SwarmMap, originReg OriginRegistrator) *RegistrationHandler {
	requestStream := make(chan handle.RequestPair, size)
	go processRequestStream(requestStream, swarmMap, originReg)
	return &RegistrationHandler{
		closed:        false,
		requestStream: requestStream,
	}
}

//AddJob adds a request to the job queue. AddJob hangs if the job queue is full
func (dh *RegistrationHandler) AddJob(request interface{}, conn handle.Conn) error {
	if dh.closed {
		return fmt.Errorf("Cannot add a job on a closed DataspaceHandler")
	}
	dh.requestStream <- handle.RequestPair{Request: request, Conn: conn}
	return nil
}

/*JobCapacity returns the max number of jobs a DatabaseHandler can hold
in its job queue*/
func (dh *RegistrationHandler) JobCapacity() int {
	return cap(dh.requestStream)
}

/*QueuedJobs returns the number of jobs currently queued in a DataspaceHandlers
job queue */
func (dh *RegistrationHandler) QueuedJobs() int {
	return len(dh.requestStream)
}

/*Close closes a DataspaceHandler for work and returns an error if the
Handler is already closed*/
func (dh *RegistrationHandler) Close() error {
	if !dh.closed {
		dh.closed = true
		close(dh.requestStream)
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
		var swarmID string
		if request.IsAdd() {
			swarmID, err = swarmMap.GetMinDataspaceSwarm()
			if err == nil {
				err = swarmMap.AddDataspace(swarmID, request.GetDataField())
			}
		} else {
			swarmID, err = swarmMap.GetSwarmID(request.GetDataField())
			if err == nil {
				err = swarmMap.RemoveDataspace(swarmID, request.GetDataField())
			}
		}
	}

	if err != nil {
		return fmt.Errorf("Failed to handle registration request in RegistrationHandler: %v", err)
	}
	return nil
}
