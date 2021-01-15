package performance

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

var HangPollTime = time.Millisecond * 50

type Addr struct {
	network string
	addr    string
}

func (a *Addr) Network() string { return a.network }
func (a *Addr) String() string  { return a.addr }

type Conn interface {
	net.Conn
	ControllerRead([]byte) (int, error)
	ControllerWrite([]byte) (int, error)
	Reset()
}

type SlowedConn struct {
	readDeadline time.Time
	closed       bool
	delay        time.Duration
	readBuf      []byte
	writeBuf     []byte
	mutex        *sync.Mutex
}

func NewSlowedConn(delay time.Duration) *SlowedConn {
	return &SlowedConn{
		readDeadline: time.Now(),
		closed:       false,
		delay:        delay,
		readBuf:      make([]byte, 0),
		writeBuf:     make([]byte, 0),
		mutex:        &sync.Mutex{},
	}
}

func (sc *SlowedConn) Read(b []byte) (int, error) {
	time.Sleep(sc.delay)

	sc.mutex.Lock()
	defer sc.mutex.Unlock()
	if sc.readDeadline == (time.Time{}) {
		if sc.closed {
			return 0, io.EOF
		}
		return 0, fmt.Errorf("Timeout Error")
	}

	for len(sc.readBuf) < len(b) {
		sc.mutex.Unlock()
		time.Sleep(HangPollTime)
		sc.mutex.Lock()
	}

	i := 0
	for ; i < len(sc.readBuf) && i < len(b); i++ {
		b[i] = sc.readBuf[i]
	}
	if i == len(sc.readBuf) {
		sc.readBuf = make([]byte, 0)
	} else {
		sc.readBuf = sc.readBuf[i:]
	}

	if i < len(b) {
		return i, fmt.Errorf("Not enough data")
	}
	return i, nil
}

func (sc *SlowedConn) ControllerRead(b []byte) (int, error) {
	sc.mutex.Lock()
	defer sc.mutex.Unlock()
	for len(sc.writeBuf) < len(b) {
		sc.mutex.Unlock()
		time.Sleep(HangPollTime)
		sc.mutex.Lock()
	}

	i := 0
	for ; i < len(b) && i < len(sc.writeBuf); i++ {
		b[i] = sc.writeBuf[i]
	}

	if i == len(sc.writeBuf) {
		sc.writeBuf = make([]byte, 0)
	} else {
		sc.writeBuf = sc.writeBuf[i:]
	}

	if i < len(b) {
		return i, fmt.Errorf("Not enough data")
	}
	return i, nil
}

func (sc *SlowedConn) Write(b []byte) (int, error) {
	time.Sleep(sc.delay)
	sc.mutex.Lock()
	defer sc.mutex.Unlock()

	sc.writeBuf = append(sc.writeBuf, b...)
	return len(b), nil
}

func (sc *SlowedConn) ControllerWrite(b []byte) (int, error) {
	sc.mutex.Lock()
	defer sc.mutex.Unlock()

	sc.readBuf = append(sc.readBuf, b...)
	return len(b), nil
}

func (sc *SlowedConn) Close() error {
	sc.closed = true
	return nil
}

func (sc *SlowedConn) LocalAddr() net.Addr {
	return &Addr{}
}

func (sc *SlowedConn) RemoteAddr() net.Addr {
	return &Addr{}
}

func (sc *SlowedConn) SetDeadline(t time.Time) error {
	sc.readDeadline = t
	return nil
}

func (sc *SlowedConn) SetReadDeadline(t time.Time) error {
	sc.readDeadline = t
	return nil
}

func (sc *SlowedConn) SetWriteDeadline(t time.Time) error {
	return nil
}
