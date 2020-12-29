package localizer

/*SwarmManager defines an object that can process a request given a dataspace
and a connection the the requester*/
type SwarmManager interface {
	AttemptToPair(conn interface{}) error
	GetID() string
}

/*SwarmMap defines an object that can return the SwarmManager assigned
to a specific dataspace*/
type SwarmMap interface {
	GetSwarm(string) (interface{}, error)
}

/*FrequencyTracker defines an object that needs to be informed when a new request
comes in in order to keep track of the frequency of requests per dataspace and swarm*/
type FrequencyTracker interface {
	IncrementFrequency(dataspace string, swarmID string)
}

/*LocalizeRequest defines an object that holds the information needed for the localizer
to properly identify where to sent the request for processing*/
type LocalizeRequest interface {
	GetDataspace() string
}
