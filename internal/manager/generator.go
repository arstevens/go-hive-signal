package manager

/*Implements mapper.SwarmManagerGenerator. Allows the creation of new
swarm managers given the provided SwarmGatewayGenerator, AgentNegotiator,
and SwarmSizeTracker objects*/
type SwarmManagerGenerator struct {
	gatewayGenerator SwarmGatewayGenerator
	negotiator       AgentNegotiator
	tracker          SwarmSizeTracker
}

//NewGenerator creates a new instance of SwarmManagerGenerator
func NewGenerator(gatewayGen SwarmGatewayGenerator, negotiate AgentNegotiator,
	tracker SwarmSizeTracker) *SwarmManagerGenerator {
	return &SwarmManagerGenerator{
		gatewayGenerator: gatewayGen,
		negotiator:       negotiate,
		tracker:          tracker,
	}
}

//New returns a new SwarmManager instance
func (sg *SwarmManagerGenerator) New(id string) interface{} {
	return New(id, sg.gatewayGenerator.New(), sg.negotiator, sg.tracker)
}
