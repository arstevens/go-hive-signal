package mapper

import "io"

type SwarmManager interface {
	io.Closer
}

type SwarmManagerGenerator interface {
	New(string) (interface{}, error)
}
