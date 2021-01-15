package route

import (
	"io"

	"github.com/arstevens/go-request/handle"
)

/* Listener defines a type that can accept new connections
to receive requests */
type Listener interface {
	io.Closer
	Accept() (handle.Conn, error)
}
