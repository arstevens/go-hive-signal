package localize

import "fmt"

/*RequestLocalizer is an implementation of the RequestHandler
interface in go-request/handle that can figure out the appropriate
swarm to pass the request on to based on the data being requested
and its IP address*/
type RequestLocalizer struct {
	requestQueue chan<- DiscoverRequest
}

/*NewRequestLocalizer returns a new *RequestLocalizer with a queue of
size queueSize and with the provided SwarmIDMap and SwarmHandlerMap*/
func NewRequestLocalizer(queueSize int, freqManager FrequencyManager,
	handlerMap SwarmMap) (*RequestLocalizer, error) {
	requestStream := make(chan DiscoverRequest, queueSize)
	go handleLocalizeRequests(requestStream, freqManager, handlerMap)
	return &RequestLocalizer{
		requestQueue: requestStream,
	}, nil
}

/*AddJob first casts the passed in interface{} into a request of
type LocalizeRequest and returns an error if request does not implement
LocalizeRequest. If the cast was successful the request will be added
to a queue to be processed*/
func (rl *RequestLocalizer) AddJob(request interface{}) error {
	localizeRequest, ok := request.(DiscoverRequest)
	if !ok {
		return fmt.Errorf("Received request not of type LocalizeRequest in RequestLocalizer")
	}
	rl.requestQueue <- localizeRequest
	return nil
}

/*JobCapacity returns the total number of jobs that the
RequestLocalizer can have queued before AddJob() hangs*/
func (rl *RequestLocalizer) JobCapacity() int {
	return cap(rl.requestQueue)
}

/*QueuedJobs returns the number of jobs currently waiting
to be processed by the localizer*/
func (rl *RequestLocalizer) QueuedJobs() int {
	return len(rl.requestQueue)
}

// Close closes the queue channel
func (rl *RequestLocalizer) Close() error {
	close(rl.requestQueue)
	return nil
}
