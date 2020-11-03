package connect

import (
	"fmt"
	"log"

	"github.com/arstevens/go-request/handle"
)

type ConnectionHandler struct {
	closed        bool
	requestStream chan<- handle.RequestPair
}

func New(size int, verifier IdentityVerifier, connector SwarmConnector) *ConnectionHandler {
	requestStream := make(chan handle.RequestPair, size)
	go processRequestStream(requestStream, verifier, connector)
	return &ConnectionHandler{
		closed:        false,
		requestStream: requestStream,
	}
}

//AddJob adds a request to the job queue. AddJob hangs if the job queue is full
func (ch *ConnectionHandler) AddJob(request interface{}, conn handle.Conn) error {
	if ch.closed {
		return fmt.Errorf("Cannot add a job on a closed ConnectionHandler")
	}
	ch.requestStream <- handle.RequestPair{Request: request, Conn: conn}
	return nil
}

/*JobCapacity returns the max number of jobs a ConnectionHandler can hold
in its job queue*/
func (ch *ConnectionHandler) JobCapacity() int {
	return cap(ch.requestStream)
}

/*QueuedJobs returns the number of jobs currently queued in a DataspaceHandlers
job queue */
func (ch *ConnectionHandler) QueuedJobs() int {
	return len(ch.requestStream)
}

/*Close closes a DataspaceHandler for work and returns an error if the
Handler is already closed*/
func (ch *ConnectionHandler) Close() error {
	if !ch.closed {
		ch.closed = true
		close(ch.requestStream)
		return nil
	}
	return fmt.Errorf("Cannot close a closed DataspaceHandler")
}

func processRequestStream(requestStream <-chan handle.RequestPair,
	verifier IdentityVerifier, connector SwarmConnector) {
	for {
		requestPair, ok := <-requestStream
		if !ok {
			return
		}
		request := requestPair.Request.(ConnectionRequest)
		conn := requestPair.Conn.(NetConn)

		err := handleConnectionRequest(request, conn, verifier, connector)
		if err != nil {
			log.Println(err)
		}
	}
}

func handleConnectionRequest(request ConnectionRequest, conn NetConn,
	verifier IdentityVerifier, connector SwarmConnector) error {
	if !verifier.Analyze(conn.GetIP(), request.GetOriginID(), request.IsLogOn()) {
		return fmt.Errorf("Identity Verification failed in ConnectionHandler")
	}

	err := connector.ProcessConnection(request.GetRequestCode(), conn)
	if err != nil {
		return fmt.Errorf("Failed to pass request to SwarmConnector in ConnectionHandler: %v", err)
	}
	return nil
}
