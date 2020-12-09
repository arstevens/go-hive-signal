package protomsg

import (
	"fmt"
	"testing"
)

func TestProtomsg(t *testing.T) {
	fmt.Printf("Testing LocalizerRequest\n")
	fmt.Printf("\tcreating new request...")
	lRequest, err := NewLocalizeRequest("/dataspace/TEST", 1)
	if err != nil {
		fmt.Printf("failed\n")
		panic(err)
	}
	fmt.Printf("success\n")

	fmt.Printf("\tunmarshaling request...")
	_, err = UnpackLocalizeRequest(lRequest, &FakeConn{})
	if err != nil {
		fmt.Printf("failed\n")
		panic(err)
	}
	fmt.Printf("success\n")

	fmt.Printf("Testing RegistrationRequest\n")
	fmt.Printf("\tcreating new request...")
	rRequest, err := NewRegistrationRequest(true, true, "DATA", 1)
	if err != nil {
		fmt.Printf("failed\n")
		panic(err)
	}
	fmt.Printf("success\n")

	fmt.Printf("\tunmarshaling request...")
	_, err = UnpackRegistrationRequest(rRequest, &FakeConn{})
	if err != nil {
		fmt.Printf("failed\n")
		panic(err)
	}
	fmt.Printf("success\n")

	fmt.Printf("Testing ConnectionRequest\n")
	fmt.Printf("\tcreating new request...")
	cRequest, err := NewConnectionRequest(true, "/swarm/TEST", "/origin/TEST", 1, 1)
	if err != nil {
		fmt.Printf("failed\n")
		panic(err)
	}
	fmt.Printf("success\n")

	fmt.Printf("\tunmarshaling request...")
	_, err = UnpackConnectionRequest(cRequest, &FakeConn{})
	if err != nil {
		fmt.Printf("failed\n")
		panic(err)
	}
	fmt.Printf("success\n")
}

type FakeConn struct{}

func (fc *FakeConn) Read([]byte) (int, error)  { return 0, nil }
func (fc *FakeConn) Write([]byte) (int, error) { return 0, nil }
func (fc *FakeConn) Close() error              { return nil }
