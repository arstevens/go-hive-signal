package manager

/*Implements mapper.SwarmManagerGenerator. Allows the creation of new
swarm managers given the provided SwarmGatewayGenerator, AgentNegotiator,
and SwarmSizeTracker objects*/
type SwarmManagerGenerator struct {
	gatewayGenerator SwarmGatewayGenerator
	negotiator       AgentNegotiator
	tracker          SwarmSizeTracker
}

func (sg *SwarmManagerGenerator) New(id string) interface{} {
	return New(id, sg.gatewayGenerator.New(), sg.negotiator, sg.tracker)
}
