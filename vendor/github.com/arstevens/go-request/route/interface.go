package route

import "github.com/arstevens/go-request/handle"

/* ReadRequest describes a function that can read a single
raw request from a Conn object */
type ReadRequest func(handle.Conn) ([]byte, error)

type UnpackRouteRequest func([]byte) (handle.Request, error)
