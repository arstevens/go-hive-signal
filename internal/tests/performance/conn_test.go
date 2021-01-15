package performance

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestSlowedConn(t *testing.T) {
	conn := NewSlowedConn(time.Millisecond * 50)

	totalWrites := 10
	for i := 0; i < totalWrites; i++ {
		msg := []byte("/msg_" + strconv.Itoa(i))
		_, err := conn.Write(msg)
		if err != nil {
			panic(err)
		}
	}

	go func() {
		time.Sleep(time.Millisecond * 500)
		conn.Write([]byte("/msg_N"))
	}()

	full := make([]byte, 0)
	for i := 0; i < totalWrites+1; i++ {
		buf := make([]byte, 6)
		_, err := conn.ControllerRead(buf)
		if err != nil {
			panic(err)
		}
		full = append(full, buf...)
	}
	fmt.Printf("%s\n", string(full))

	for i := 0; i < totalWrites; i++ {
		_, err := conn.ControllerWrite([]byte("/msg_" + strconv.Itoa(i)))
		if err != nil {
			panic(err)
		}
	}

	full = make([]byte, 0)
	for i := 0; i < totalWrites; i++ {
		buf := make([]byte, 6)
		_, err := conn.Read(buf)
		if err != nil {
			panic(err)
		}
		full = append(full, buf...)
	}
	fmt.Printf("%s\n", string(full))
}
