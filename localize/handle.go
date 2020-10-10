package localize

import (
	"fmt"
	"log"
)

// attempts to handle all requests passed to requestStream until it is clsoed
func handleLocalizeRequests(requestStream <-chan LocalizeRequest, idMap SwarmIDMap, handlerMap SwarmHandlerMap) {
	for {
		request, ok := <-requestStream
		if !ok {
			log.Println("Localize Request Stream closed for handleLocalizeRequests(). Returning.")
			return
		}
		err := sendToSwarmHandler(request, idMap, handlerMap)
		if err != nil {
			log.Printf("Failed to process request in RequestLocalizer: %v\n", err)
		}
	}
}

// performs the database lookups to find the handler and attempts to pass the job along
func sendToSwarmHandler(request LocalizeRequest, idMap SwarmIDMap, handlerMap SwarmHandlerMap) error {
	swarmID, err := idMap.GetSwarmID(request.GetDataID(), request.GetIPAddress())
	if err != nil {
		return fmt.Errorf("Failed to retrieve SwarmID from request info: %v", err)
	}
	swarmHandler, err := handlerMap.GetSwarmHandler(swarmID)
	if err != nil {
		return fmt.Errorf("Failed to retrieve RequestHandler from SwarmID: %v", err)
	}
	err = swarmHandler.AddJob(request)
	if err != nil {
		return fmt.Errorf("Failed to add request to RequestHandler: %v", err)
	}
	return nil
}
