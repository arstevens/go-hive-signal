package protomsg

type PBRouteWrapper struct {
	request *RouterWrapper
}

type PBLocalizeRequest struct {
	request *LocalizeRequest
}

func (lr *PBLocalizeRequest) GetDataspace() string {
	return lr.request.GetDataspace()
}

type PBRegistrationRequest struct {
	request *RegistrationRequest
}

func (rr *PBRegistrationRequest) IsAdd() bool {
	return rr.request.GetIsAdd()
}

func (rr *PBRegistrationRequest) IsOrigin() bool {
	return rr.request.GetIsOrigin()
}

func (rr *PBRegistrationRequest) GetDataField() string {
	return rr.request.GetDatafield()
}

type PBConnectionRequest struct {
	request *ConnectionRequest
}

func (cr *PBConnectionRequest) IsLogOn() bool {
	return cr.request.GetIsLogOn()
}

func (cr *PBConnectionRequest) GetSwarmID() string {
	return cr.request.GetSwarmID()
}

func (cr *PBConnectionRequest) GetOriginID() string {
	return cr.request.GetOriginID()
}

type PBNegotiateMessage struct {
	msg *NegotiateMessage
}

func (nm *PBNegotiateMessage) IsAccepted() bool {
	return nm.msg.GetIsAccepted()
}
