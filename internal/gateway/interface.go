package gateway

import "github.com/arstevens/go-hive-signal/internal/manager"

type Conn interface {
	manager.Conn
	IsClosed() bool
}
