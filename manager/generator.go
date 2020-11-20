package manager

/*SwarmManagerGenerator generates new SwarmManagers with the stored
'gateway' and 'negotiate' members*/
type SwarmManagerGenerator struct {
	gateway   SwarmGateway
	negotiate AgentNegotiator
}

//NewGenerator creates a new SwarmManagerGenerator
func NewGenerator(gateway SwarmGateway, negotiate AgentNegotiator) *SwarmManagerGenerator {
	return &SwarmManagerGenerator{
		gateway:   gateway,
		negotiate: negotiate,
	}
}

//New creates a new SwarmManager with 'swarmID'
func (sg *SwarmManagerGenerator) New(swarmID string) (interface{}, error) {
	return New(swarmID, sg.gateway, sg.negotiate), nil
}
