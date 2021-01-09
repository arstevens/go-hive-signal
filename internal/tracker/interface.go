package tracker

/*StorageEngine outlines an interface for an object
that can take a Datapoint and store it as well as
return some data back*/
type StorageEngine interface {
	AddDatapoint(interface{})
	GetData() interface{}
}

/*StorageEngineGenerator describes an object that
can return a new instance of a StorageEngine*/
type StorageEngineGenerator interface {
	New() StorageEngine
}
