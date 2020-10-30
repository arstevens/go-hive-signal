package dataspace

/*ValidDataspaceStore describes a store that allows for the
addition and removal of valid dataspaces for the signaling server*/
type ValidDataspaceStore interface {
	AddDataspace(string)
	RemoveDataspace(string)
}

/*SwarmMap describes an object that can map a dataspace to a
specific swarm as well as having the ability to return the swarm
with the least number of dataspaces*/
type SwarmMap interface {
	GetMinDataspaceSwarm() string
	AddDataspace(swarmID string, dataspace string)
}
