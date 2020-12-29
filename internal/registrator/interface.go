package registrator

/*OriginRegistrator describes a datastore that keeps track of
which valid 'points of origin' for endpoints to come from */
type OriginRegistrator interface {
	AddOrigin(string) error
	RemoveOrigin(string) error
}

/*SwarmMap describes an object that can map a dataspace to a
specific swarm as well as having the ability to return the swarm
with the least number of dataspaces*/
type SwarmMap interface {
	AddSwarm(dataspace string) error
	RemoveSwarm(dataspace string) error
}

/* DataspaceRequest describes an object that represents a request
that would be passed to the DataspaceHandler. If request.IsOrigin() is
true then GetDataField() will return the origin name. Else it will return
a dataspace name */
type RegistrationRequest interface {
	IsAdd() bool
	IsOrigin() bool
	GetDataField() string
}
