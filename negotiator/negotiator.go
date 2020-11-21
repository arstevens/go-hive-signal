package negotiator

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/arstevens/go-hive-signal/manager"
)

var RoundtripLimit = 5
var HeaderEndian = binary.LittleEndian
var UnmarshalMessage UnmarshalNegotiateMessage = nil

var readErrorString = "Failed to read message in RoundtripLimitedNegotiate(): %v"
var writeErrorString = "Failed to write message in RoundtripLimitedNegotiate(): %v"

func RoundtripLimitedNegotiate(offerer manager.Conn, acceptor manager.Conn) error {
	defer offerer.Close()
	defer acceptor.Close()
	for i := 0; i < RoundtripLimit; i++ {
		rawOffer, err := readMessageFromConn(offerer)
		if err != nil {
			return fmt.Errorf(readErrorString, err)
		}
		_, err = acceptor.Write(rawOffer)
		if err != nil {
			return fmt.Errorf(writeErrorString, err)
		}

		rawResponse, err := readMessageFromConn(acceptor)
		if err != nil {
			return fmt.Errorf(readErrorString, err)
		}
		ifaceMsg, err := UnmarshalMessage(rawResponse)
		if err != nil {
			return fmt.Errorf("Failed to unmarshal response in RoundtripLimitedNegotitate(): %v", err)
		}
		message := ifaceMsg.(NegotiateMessage)

		_, err = offerer.Write(rawResponse)
		if err != nil {
			return fmt.Errorf(writeErrorString, err)
		}
		if message.IsAccepted() {
			return nil
		}
	}
	return fmt.Errorf("Roundtrip limit reached without consensus in RountripLimitedNegotiate()")
}

func readMessageFromConn(conn io.Reader) ([]byte, error) {
	var size uint64
	err := binary.Read(conn, HeaderEndian, &size)
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
