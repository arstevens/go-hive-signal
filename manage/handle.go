package manage

import (
	"fmt"
	"log"
)

// attempts to process a stream of ConnectionRequests
func handleConnectionRequests(requestStream <-chan ConnectionRequest,
	verifier IdentityVerifier, tracker MemberTracker) {
	for {
		request, ok := <-requestStream
		if !ok {
			log.Println("Connection Request Stream closed for handleConnectionRequests(). Returning.")
			return
		}
		err := verifyAndUpdate(request, verifier, tracker)
		if err != nil {
			log.Printf("Failed to process request in MemberManager: %v\n", err)
		}
	}
}

// perform the steps required to deal with a ConnectionRequest
func verifyAndUpdate(request ConnectionRequest, verifier IdentityVerifier,
	tracker MemberTracker) error {
	origin := request.GetOrigin()
	ipAddr := request.GetIPAddress()

	err := verifier.Verify(origin, ipAddr)
	if err != nil {
		return fmt.Errorf("Failed to verify request using IdentityVerifier: %v", err)
	}
	err = tracker.ModifyTrackingData(origin, ipAddr, request.GetActivityStatus())
	if err != nil {
		return fmt.Errorf("Failed to modify tracking data using MemberTracker: %v", err)
	}
	return nil
}
