package negotiator

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/arstevens/go-hive-signal/internal/manager"
)

var RoundtripLimit = 5
var UnmarshalMessage UnmarshalNegotiateMessage = nil

var readErrorString = "Failed to read message in RoundtripLimitedNegotiate(): %v"
var writeErrorString = "Failed to write message in RoundtripLimitedNegotiate(): %v"

func RoundtripLimitedNegotiate(offerer manager.Conn, acceptor manager.Conn) error {
	for i := 0; i < RoundtripLimit; i++ {
		rawOffer, err := readMessageFromWire(offerer)
		if err != nil {
			return fmt.Errorf(readErrorString, err)
		}
		err = writeMessageToWire(acceptor, rawOffer)
		if err != nil {
			return fmt.Errorf(writeErrorString, err)
		}

		rawResponse, err := readMessageFromWire(acceptor)
		if err != nil {
			return fmt.Errorf(readErrorString, err)
		}
		ifaceMsg, err := UnmarshalMessage(rawResponse)
		if err != nil {
			return fmt.Errorf("Failed to unmarshal response in RoundtripLimitedNegotitate(): %v", err)
		}
		message := ifaceMsg.(NegotiateMessage)

		err = writeMessageToWire(offerer, rawResponse)
		if err != nil {
			return fmt.Errorf(writeErrorString, err)
		}
		if message.IsAccepted() {
			return nil
		}
	}
	return fmt.Errorf("Roundtrip limit reached without consensus in RountripLimitedNegotiate()")
}

func writeMessageToWire(conn io.Writer, msg []byte) error {
	var size int64
	size = int64(len(msg))
	err := binary.Write(conn, binary.BigEndian, size)
	if err != nil {
		return fmt.Errorf("Failed to write header to connection: %v", err)
	}

	_, err = conn.Write(msg)
	if err != nil {
		return fmt.Errorf("Failed to write message to connection: %v", err)
	}
	return nil
}

func readMessageFromWire(conn io.Reader) ([]byte, error) {
	var size int64
	err := binary.Read(conn, binary.BigEndian, &size)
	if err != nil {
		return nil, fmt.Errorf("Failed to read message header from conn: %v", err)
	}

	buf := make([]byte, size)
	_, err = io.ReadFull(conn, buf)
	if err != nil {
		return nil, fmt.Errorf("Failed to read message from conn: %v", err)
	}
	return buf, nil
}
