package manage

/*MemberAllocator describes an object that stores information
regarding the current state of a p2p endpoint and can allocate
a new endpoint to a swarm*/
type MemberAllocator interface {
	AllocateToSwarm(request interface{}) error
	AllocateToJob(request interface{}) error
	RemoveFromSwarm(request interface{}) error
}
