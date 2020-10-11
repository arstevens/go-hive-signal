package localize

import (
	"fmt"
	"log"
)

// attempts to handle all requests passed to requestStream until it is clsoed
func handleLocalizeRequests(requestStream <-chan DiscoverRequest, freqManager FrequencyManager,
	handlerMap SwarmMap) {
	for {
		request, ok := <-requestStream
		if !ok {
			log.Println("Localize Request Stream closed for handleLocalizeRequests(). Returning.")
			return
		}
		freqManager.IncrementFrequency(request.GetDataID())
		err := sendToSwarmHandler(request, handlerMap)
		if err != nil {
			log.Printf("Failed to process request in RequestLocalizer: %v\n", err)
		}
	}
}

// performs the database lookups to find the handler and attempts to pass the job along
func sendToSwarmHandler(request DiscoverRequest, handlerMap SwarmMap) error {
	swarmHandler, err := handlerMap.GetSwarmHandler(request.GetDataID())
	if err != nil {
		return fmt.Errorf("Failed to retrieve RequestHandler from SwarmID: %v", err)
	}
	err = swarmHandler.AddJob(request)
	if err != nil {
		return fmt.Errorf("Failed to add request to RequestHandler: %v", err)
	}
	return nil
}
