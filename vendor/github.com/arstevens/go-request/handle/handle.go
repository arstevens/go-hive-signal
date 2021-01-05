package handle

import "io"

/* Conn describes an object that can be used to read
request from */
type Conn interface {
	io.ReadWriteCloser
}

/* RequestHandler describes an object that can handle a
stream of requests. RequestHandler is responsible for closing
the Conn for each job it receives */
type RequestHandler interface {
	AddJob(interface{}, Conn) error
	JobCapacity() int
	QueuedJobs() int
	Close() error
}

/* RequestHandlerGenerator describes an object that
generates new objects that conforms to the RequestHandler interface */
type RequestHandlerGenerator interface {
	NewHandler() RequestHandler
	HandlerCapacity() int
}

/* Defines a Request object with only one method to
retrieve an integer identifying the type of the request */
type Request interface {
	GetType() int32
	GetRequest() []byte
}

/* RequestPair is a datatype composed of the received request
and a connection to the party that sent the request */
type RequestPair struct {
	Request interface{}
	Conn    Conn
}

/* Defines a function that can take a sequence of bytes
and attempt to unpack it into an object usable by a RequestHandler */
type UnpackRequest func([]byte) (interface{}, error)
