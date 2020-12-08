package manager

import "fmt"

/*SwarmManager is an object that can be used to connect new requesters
to a peer-to-peer swarm*/
type SwarmManager struct {
	gateway   SwarmGateway
	negotiate AgentNegotiator
	closed    bool
	id        string
}

//New creates a new SwarmManager
func New(swarmID string, gateway SwarmGateway, negotiate AgentNegotiator) *SwarmManager {
	return &SwarmManager{
		gateway:   gateway,
		negotiate: negotiate,
		closed:    false,
		id:        swarmID,
	}
}

//AttemptToPair attempts to pair 'conn' with someone from the swarm
func (sm *SwarmManager) AttemptToPair(conn interface{}) error {
	if sm.closed {
		return fmt.Errorf("Failed to attempt pair in SwarmManager.AttemptToPair(). Object closed")
	}

	acceptorConn, ok := conn.(Conn)
	if !ok {
		return fmt.Errorf("Failed to pair in SwarmManager.AttemptToPair(). 'conn' does not conform to Conn interface")
	}
	offerer, err := sm.gateway.GetEndpoint(sm.id)
	if err != nil {
		return fmt.Errorf("Failed to pair in SwarmManager.AttemptToPair(): %v", err)
	}
	offererConn, ok := offerer.(Conn)
	if !ok {
		return fmt.Errorf("Failed to pair in SwarmManager.AttemptToPair(). offerer 'conn' does not conform to Conn interface")
	}

	err = sm.negotiate(offererConn, acceptorConn)
	if err != nil {
		return fmt.Errorf("Failed to pair in SwarmManager.AttemptToPair(): %v", err)
	}
	return nil
}

//GetID returns the swarm ID associated with this manager
func (sm *SwarmManager) GetID() string {
	return sm.id
}

//Close closes the SwarmManager for use
func (sm *SwarmManager) Close() error {
	if sm.closed {
		return fmt.Errorf("Failed to close in SwarmManager.Close() Already closed")
	}
	sm.closed = true
	return nil
}
