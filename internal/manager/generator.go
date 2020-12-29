package manager

type SwarmManagerGenerator struct {
	gateway    SwarmGateway
	negotiator AgentNegotiator
	tracker    SwarmSizeTracker
}

func (sg *SwarmManagerGenerator) New(id string) interface{} {
	return New(id, sg.gateway, sg.negotiator, sg.tracker)
}
