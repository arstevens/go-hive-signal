package route

import (
	"github.com/arstevens/go-request/handle"
)

/* UnpackAndRoute begins a pipeline for the listening, interpreting and routing
of requests from clients */
func UnpackAndRoute(listener Listener, done <-chan struct{}, handlers map[int32]handle.RequestHandler,
	unpack UnpackRouteRequest, unpackers map[int32]handle.UnpackRequest, read ReadRequest) {
	identifyStream := make(chan directPair)
	pipelineDone := make(chan struct{})
	defer close(pipelineDone)

	go identifyAndRoute(identifyStream, handlers)
	go listenAndUnmarshal(listener, unpack, read, unpackers, pipelineDone, identifyStream)
	<-done
}
