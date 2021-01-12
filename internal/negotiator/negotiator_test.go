package negotiator

import (
	"encoding/binary"
	"fmt"
	"io"
	"testing"
)

func TestRountripLimitNegotiate(t *testing.T) {
	fmt.Printf("----------------------\nROUNDTRIP LIMIT NEGOTIATE TEST\n----------------------\n")
	UnmarshalMessage = unmarshal

	BufSize := 1000
	offerer := FakeConn{buf: make([]byte, BufSize), head: 0, tail: 0}
	acceptor := FakeConn{buf: make([]byte, BufSize), head: 0, tail: 0}

	numberOfTrips := 5
	writeNOffers(&offerer, numberOfTrips)
	writeNResponses(&acceptor, numberOfTrips)

	err := RoundtripLimitedNegotiate(&offerer, &acceptor)
	fmt.Printf("Status: %v\n", err)
}

func writeNResponses(conn io.Writer, n int) {
	resp1 := []byte("false")
	resp2 := []byte("true")

	var size uint64
	for i := 0; i < n; i++ {
		var resp []byte
		if i == n-1 {
			resp = resp2
			size = uint64(len(resp2))
		} else {
			resp = resp1
			size = uint64(len(resp1))
		}

		err := binary.Write(conn, binary.BigEndian, size)
		if err != nil {
			panic(err)
		}
		conn.Write(resp)
	}
}

func writeNOffers(conn io.Writer, n int) {
	msg := []byte("offer")
	var size uint64
	size = uint64(len(msg))
	for i := 0; i < n; i++ {
		err := binary.Write(conn, binary.BigEndian, size)
		if err != nil {
			panic(err)
		}
		conn.Write(msg)
	}
}

func TestMessageRead(t *testing.T) {
	fmt.Printf("----------------------\nMESSAGE READ TEST\n----------------------\n")
	conn := FakeConn{buf: make([]byte, 100), head: 0, tail: 0}
	msg := []byte("this is the message")
	var size uint64
	size = uint64(len(msg))
	err := binary.Write(&conn, binary.BigEndian, size)
	if err != nil {
		panic(err)
	}
	conn.Write(msg)

	resp, err := readMessageFromWire(&conn)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Read: %s\n", string(resp))
}

func TestFakeConn(t *testing.T) {
	fmt.Printf("----------------------\nFAKECONN TEST\n----------------------\n")
	conn := FakeConn{buf: make([]byte, 100), head: 0, tail: 0}

	out1 := []byte("this is a string")
	conn.Write(out1)

	out2 := []byte("this is another string")
	conn.Write(out2)

	buf1 := make([]byte, len(out1))
	conn.Read(buf1)
	fmt.Printf("Read: %s\n", string(buf1))

	buf2 := make([]byte, len(out2))
	conn.Read(buf2)
	fmt.Printf("Read: %s\n", string(buf2))
}

func unmarshal(b []byte) (interface{}, error) {
	if string(b) == "true" {
		return &message{isAccepted: true}, nil
	}
	return &message{isAccepted: false}, nil
}

type message struct {
	isAccepted bool
}

func (m *message) IsAccepted() bool { return m.isAccepted }

type FakeConn struct {
	buf  []byte
	head int
	tail int
}

func (fc *FakeConn) GetAddress() string { return "" }
func (fc *FakeConn) IsClosed() bool     { return false }

func (fc *FakeConn) Read(b []byte) (int, error) {
	i := 0
	for ; i < len(b) && i+fc.head < fc.tail; i++ {
		b[i] = fc.buf[fc.head+i]
	}
	fc.head += i
	return i, nil
}

func (fc *FakeConn) Write(b []byte) (int, error) {
	i := 0
	for ; i < len(b) && fc.tail+i < len(fc.buf); i++ {
		fc.buf[fc.tail+i] = b[i]
	}
	fc.tail += i
	return i, nil
}

func (fc *FakeConn) Close() error {
	return nil
}
