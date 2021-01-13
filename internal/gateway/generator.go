package gateway

import "github.com/arstevens/go-hive-signal/internal/manager"

/*SwarmGatewayGenerator implements manager.SwarmGatewayGenerator. It
defines an object that can create new SwarmGateways with the provided
activeSize and inactiveSize*/
type SwarmGatewayGenerator struct {
	activeSize int
}

//NewGenerator creates a new SwarmGatewayGenerator instance
func NewGenerator(activeSize int) *SwarmGatewayGenerator {
	return &SwarmGatewayGenerator{activeSize: activeSize}
}

//New creates a new SwarmGateway instance
func (sg *SwarmGatewayGenerator) New() manager.SwarmGateway {
	return New(sg.activeSize)
}
