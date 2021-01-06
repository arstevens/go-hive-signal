package route

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"

	"github.com/arstevens/go-request/handle"
)

/* listenAndUnmarshal accepts any connections from listener and attempts
to read and deserialize a request */
func listenAndUnmarshal(listener Listener, unpacker UnpackRouteRequest, reader ReadRequest,
	unpackerMap map[int32]handle.UnpackRequest, done <-chan struct{}, outStream chan<- directPair) {
	defer close(outStream)
	defer listener.Close()

	requestChan := make(chan directPair)
	go receiveRequests(listener, unpacker, reader, unpackerMap, requestChan)
	for {
		select {
		case request, ok := <-requestChan:
			if !ok {
				return
			}
			outStream <- request
		case <-done:
			return
		}
	}
}

/* receiveRequests accepts all connections on the listener and
attempts to deserialize them into handle.Request objects. It then passes
these objects through the returnStream channel */
func receiveRequests(listener Listener, routeUnpacker UnpackRouteRequest, reader ReadRequest,
	unpackers map[int32]handle.UnpackRequest, returnStream chan<- directPair) {
	defer close(returnStream)
	for {
		conn, err := listener.Accept()
		/* An error will be returned if the listener was closed.
		this allows gracefully stopping */
		if err != nil {
			log.Println(err)
			return
		}

		request, err := readAndUnpackRequest(conn, reader, routeUnpacker)
		if err != nil {
			conn.Close()
			log.Println(err)
			continue
		}
		requestType := request.GetType()
		unpacker, ok := unpackers[requestType]
		if !ok {
			conn.Close()
			log.Println(err)
			continue
		}
		requestIface, err := unpacker(request.GetRequest())
		if err != nil {
			conn.Close()
			log.Printf("Failed to unpack request in listenAndUnmarshal: %v\n", err)
			continue
		}

		pair := handle.RequestPair{requestIface, conn}
		returnStream <- directPair{pair, requestType}
	}
}

// Consolidates the reading of a request
func readAndUnpackRequest(conn handle.Conn, reader ReadRequest, unpacker UnpackRouteRequest) (handle.Request, error) {
	rawRequest, err := reader(conn)
	if err != nil {
		return nil, err
	}
	return unpacker(rawRequest)
}

/* ReadRequestFromConn performs the actual reading from
a single connection */
func ReadRequestFromNetConn(conn handle.Conn) ([]byte, error) {
	var packetSize int32
	err := binary.Read(conn, binary.BigEndian, &packetSize)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("prefix read error in readRequestFromConn(): %v", err)
	}
	requestPacket := make([]byte, packetSize)
	_, err = io.ReadFull(conn, requestPacket)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("packet read error in readRequestFromConn(): %v", err)
	}
	return requestPacket, nil
}
