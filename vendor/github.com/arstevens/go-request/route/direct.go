package route

import (
	"errors"
	"log"

	"github.com/arstevens/go-request/handle"
)

type directPair struct {
	Pair        handle.RequestPair
	RequestType int32
}

/* UnknownRequestErr is an error indicating receiving a request
that could not be mapped to a handler */
var UnknownRequestErr = errors.New("Unknown request type code")

// identifyAndRoute takes a request and sends it to the proper subcomponent
func identifyAndRoute(requestStream <-chan directPair, handlers map[int32]handle.RequestHandler) {
	for {
		directPair, ok := <-requestStream
		if !ok {
			log.Printf("Request Stream has been closed\n")
			return
		}

		requestPair := directPair.Pair
		handler, ok := handlers[directPair.RequestType]
		if !ok {
			requestPair.Conn.Close()
			log.Println(UnknownRequestErr)
			continue
		}
		err := handler.AddJob(requestPair.Request, requestPair.Conn)
		if err != nil {
			log.Println(err)
		}
	}
}
