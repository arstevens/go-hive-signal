package localizer

import (
	"fmt"
	"log"

	"github.com/arstevens/go-request/handle"
)

//RequestLocalizer routes requests to their right swarm
type RequestLocalizer struct {
	closed        bool
	requestStream chan<- handle.RequestPair
}

//New creates a new instance of RequestLocalizer with a job queue of capacity 'size'
func New(size int, managers SwarmMap, tracker FrequencyTracker) *RequestLocalizer {
	requestStream := make(chan handle.RequestPair, size)
	go processRequestStream(requestStream, managers, tracker)
	return &RequestLocalizer{
		closed:        false,
		requestStream: requestStream,
	}
}

//AddJob adds a request and connection to the queue for processing
func (rl *RequestLocalizer) AddJob(request interface{}, conn handle.Conn) error {
	if rl.closed {
		return fmt.Errorf("Cannot add a job on a closed RequestLocalizer")
	}
	rl.requestStream <- handle.RequestPair{Request: request, Conn: conn}
	return nil
}

//JobCapacity returns the max amount of jobs that can be queued at once
func (rl *RequestLocalizer) JobCapacity() int {
	return cap(rl.requestStream)
}

//QueuedJobs returns the number of jobs currently queued
func (rl *RequestLocalizer) QueuedJobs() int {
	return len(rl.requestStream)
}

//Close closes the RequestLocalizer
func (rl *RequestLocalizer) Close() error {
	if !rl.closed {
		close(rl.requestStream)
		rl.closed = true
		return nil
	}
	return fmt.Errorf("Cannot close a closed RequestLocalizer")
}

func processRequestStream(requestStream <-chan handle.RequestPair, managers SwarmMap, tracker FrequencyTracker) {
	for {
		requestPair, ok := <-requestStream
		if !ok {
			return
		}
		localizeRequest := requestPair.Request.(LocalizeRequest)

		err := handleLocalizeRequest(localizeRequest.GetDataspace(), requestPair.Conn, managers, tracker)
		if err != nil {
			log.Println(err)
		}
	}
}

func handleLocalizeRequest(dataspace string, conn handle.Conn, managers SwarmMap, tracker FrequencyTracker) error {
	swarmManagerObj, err := managers.GetSwarm(dataspace)
	if err != nil {
		return fmt.Errorf("Failed to get SwarmManager from SwarmMap in RequestLocalizer: %v", err)
	}
	swarmManager := swarmManagerObj.(SwarmManager)
	err = swarmManager.AttemptToPair(conn)
	if err != nil {
		return fmt.Errorf("Failed to pair to swarm in RequestLocalizer: %v", err)
	}
	tracker.IncrementFrequency(dataspace, swarmManager.GetID())
	return nil
}
