package manage

import "fmt"

/*MemberManager is an object that implements handle.RequestHandler
for connecting p2p endpoints to the signaling server*/
type MemberManager struct {
	requestQueue chan<- ConnectionRequest
}

/*NewMemberManager creates a new instance of MemberManager with the
provided verifier and tracker and a queue of size queueSize*/
func NewMemberManager(queueSize int, verifier IdentityVerifier,
	allocator MemberAllocator) (*MemberManager, error) {
	requestStream := make(chan ConnectionRequest, queueSize)
	go handleConnectionRequests(requestStream, verifier, allocator)
	return &MemberManager{
		requestQueue: requestStream,
	}, nil
}

/*AddJob attempts to cast the passed in interface{} to type
ConnectionRequest then passes the object to the request
stream processor. An error is returned if the passed in interface{}
does not implement ConnectionRequest*/
func (mm *MemberManager) AddJob(request interface{}) error {
	connRequest, ok := request.(ConnectionRequest)
	if !ok {
		return fmt.Errorf("Received request not of type ConnectionRequest in MemberManager")
	}
	mm.requestQueue <- connRequest
	return nil
}

// JobCapacity returns the max number of jobs MemberManager can queue
func (mm *MemberManager) JobCapacity() int {
	return cap(mm.requestQueue)
}

// QueuedJobs returns the number of currently queued jobs
func (mm *MemberManager) QueuedJobs() int {
	return len(mm.requestQueue)
}

// Close closes the stream processor that handles requests
func (mm *MemberManager) Close() error {
	close(mm.requestQueue)
	return nil
}
