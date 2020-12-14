package gateway

import "github.com/arstevens/go-hive-signal/internal/manager"

type Conn interface {
	GetAddress() string
	IsClosed() bool
	manager.Conn
}
