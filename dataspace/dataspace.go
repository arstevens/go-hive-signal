package dataspace

import (
	"fmt"

	"github.com/arstevens/go-request/handle"
)

/*DataspaceHandler is an object that can handle requests that add a dataspace to
the set of all supported dataspaces for a signaling server*/
type DataspaceHandler struct {
	closed       bool
	requestQueue chan<- handle.RequestPair
}

//AddJob adds a request to the job queue. AddJob hangs if the job queue is full
func (dh *DataspaceHandler) AddJob(request interface{}, conn handle.Conn) error {
	if dh.closed {
		return fmt.Errorf("Cannot add a job on a closed DataspaceHandler")
	}
	dh.requestQueue <- handle.RequestPair{Request: request, Conn: conn}
	return nil
}

/*JobCapacity returns the max number of jobs a DatabaseHandler can hold
in its job queue*/
func (dh *DataspaceHandler) JobCapacity() int {
	return cap(dh.requestQueue)
}

/*QueuedJobs returns the number of jobs currently queued in a DataspaceHandlers
job queue */
func (dh *DataspaceHandler) QueuedJobs() int {
	return len(dh.requestQueue)
}

/*Close closes a DataspaceHandler for work and returns an error if the
Handler is already closed*/
func (dh *DataspaceHandler) Close() error {
	if !dh.closed {
		dh.closed = true
		close(dh.requestQueue)
		return nil
	}
	return fmt.Errorf("Cannot close a closed DataspaceHandler")
}
