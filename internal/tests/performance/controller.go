package performance

import (
	"encoding/binary"
	"io"

	"github.com/arstevens/go-hive-signal/pkg/protomsg"
)

var RandomMsgField []byte = []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

func RequesterNegotiateRoundtrip(conn Conn, succeed bool) {
	raw, err := readFromWire(conn)
	if err != nil {
		panic(err)
	}

	iMsg, err := protomsg.UnmarshalNegotiateMessage(raw)
	if err != nil {
		panic(err)
	}
	if _, ok := iMsg.(*protomsg.PBNegotiateMessage); !ok {
		panic(err)
	}

	var resp []byte
	if succeed {
		resp, err = protomsg.NewNegotiateMessage(true, RandomMsgField)
		if err != nil {
			panic(err)
		}
	} else {
		resp, err = protomsg.NewNegotiateMessage(false, RandomMsgField)
		if err != nil {
			panic(err)
		}
	}

	err = writeToWire(resp, conn)
	if err != nil {
		panic(err)
	}
}

func writeToWire(b []byte, w io.Writer) error {
	var size int32 = int32(len(b))
	err := binary.Write(w, binary.BigEndian, size)
	if err != nil {
		return err
	}
	_, err = w.Write(b)
	return err
}

func readFromWire(r io.Reader) ([]byte, error) {
	var size int32
	err := binary.Read(r, binary.BigEndian, &size)
	if err != nil {
		panic(err)
	}

	buf := make([]byte, size)
	_, err = r.Read(buf)
	return buf, err
}
