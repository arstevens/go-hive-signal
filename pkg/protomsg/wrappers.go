package protomsg

import (
	"fmt"

	"github.com/arstevens/go-request/handle"
	"google.golang.org/protobuf/proto"
)

type PBLocalizeRequest struct {
	request *LocalizeRequest
}

func (lr *PBLocalizeRequest) GetDataspace() string {
	return lr.request.GetDataspace()
}

func (lr *PBLocalizeRequest) GetType() int {
	return int(lr.request.GetType())
}

func UnpackLocalizeRequest(raw []byte, conn handle.Conn) (handle.Request, error) {
	var request LocalizeRequest
	err := proto.Unmarshal(raw, &request)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal in UnpackLocalizeRequest(): %v", err)
	}
	return &PBLocalizeRequest{request: &request}, nil
}

type PBRegistrationRequest struct {
	request *RegistrationRequest
}

func (rr *PBRegistrationRequest) GetType() int {
	return int(rr.request.GetType())
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

func UnpackRegistrationRequest(raw []byte, conn handle.Conn) (handle.Request, error) {
	var request RegistrationRequest
	err := proto.Unmarshal(raw, &request)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal in UnpackRegistrationRequest(): %v", err)
	}
	return &PBRegistrationRequest{request: &request}, nil
}

type PBConnectionRequest struct {
	request *ConnectionRequest
}

func (cr *PBConnectionRequest) GetType() int {
	return int(cr.request.GetType())
}

func (cr *PBConnectionRequest) IsLogOn() bool {
	return cr.request.GetIsLogOn()
}

func (cr *PBConnectionRequest) GetRequestCode() int {
	return int(cr.request.GetRequestCode())
}

func (cr *PBConnectionRequest) GetSwarmID() string {
	return cr.request.GetSwarmID()
}

func (cr *PBConnectionRequest) GetOriginID() string {
	return cr.request.GetOriginID()
}

func UnpackConnectionRequest(raw []byte, conn handle.Conn) (handle.Request, error) {
	var request ConnectionRequest
	err := proto.Unmarshal(raw, &request)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal in UnpackConnectionRequest(): %v", err)
	}
	return &PBConnectionRequest{request: &request}, nil
}

type PBNegotiateMessage struct {
	msg *NegotiateMessage
}

func (nm *PBNegotiateMessage) IsAccepted() bool {
	return nm.msg.GetIsAccepted()
}

func UnmarshalNegotiateMessage(raw []byte) (interface{}, error) {
	var msg NegotiateMessage
	err := proto.Unmarshal(raw, &msg)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal in UnpackConnectionRequest(): %v", err)
	}
	return &PBNegotiateMessage{msg: &msg}, nil
}
