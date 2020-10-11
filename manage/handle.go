package manage

import (
	"fmt"
	"log"
)

// attempts to process a stream of ConnectionRequests
func handleConnectionRequests(requestStream <-chan ConnectionRequest,
	verifier IdentityVerifier, allocator MemberAllocator) {
	for {
		request, ok := <-requestStream
		if !ok {
			log.Println("Connection Request Stream closed for handleConnectionRequests(). Returning.")
			return
		}
		err := processSingleConnectionRequest(request, verifier, allocator)
		if err != nil {
			log.Printf("Failed to process request in MemberManager: %v\n", err)
		}
	}
}

// perform the steps required to deal with a ConnectionRequest
func processSingleConnectionRequest(request ConnectionRequest, verifier IdentityVerifier,
	allocator MemberAllocator) error {
	ipAddr := request.GetIPAddress()
	origin := request.GetOrigin()

	err := verifier.Verify(origin, ipAddr)
	if err != nil {
		return fmt.Errorf("Failed to verify request using IdentityVerifier: %v", err)
	}

	if request.GetIsLeaving() {
		err = allocator.RemoveFromSwarm(request)
	} else if request.GetIsNew() {
		err = allocator.AllocateToSwarm(request)
	} else {
		err = allocator.AllocateToJob(request)
	}
	if err != nil {
		return fmt.Errorf("Failed to allocate to request using MemberAllocator: %v", err)
	}
	return nil
}
