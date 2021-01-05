package protomsg

import (
	"fmt"

	"github.com/arstevens/go-request/handle"
	"google.golang.org/protobuf/proto"
)

func NewRouteWrapper(routeCode int32, rawRequest []byte) ([]byte, error) {
	request := RouterWrapper{Type: routeCode, Request: rawRequest}
	raw, err := proto.Marshal(&request)
	if err != nil {
		return nil, fmt.Errorf("Failed to create RouteWrapper in NewRouteWrapper(): %v", err)
	}
	return raw, nil
}

func UnpackRouteWrapper(raw []byte) (handle.Request, error) {
	var request RouterWrapper
	err := proto.Unmarshal(raw, &request)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal in UnpackRouteWrapper(): %v", err)
	}
	return &request, nil
}

func NewLocalizeRequest(dataspace string) ([]byte, error) {
	request := LocalizeRequest{Dataspace: dataspace}
	raw, err := proto.Marshal(&request)
	if err != nil {
		return nil, fmt.Errorf("Failed to create LocalizeRequest in NewLocalizeRequest(): %v", err)
	}
	return raw, nil
}

func UnpackLocalizeRequest(raw []byte) (interface{}, error) {
	var request LocalizeRequest
	err := proto.Unmarshal(raw, &request)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal in UnpackLocalizeRequest(): %v", err)
	}
	return &PBLocalizeRequest{request: &request}, nil
}

func NewRegistrationRequest(isAdd bool, isOrigin bool, datafield string) ([]byte, error) {
	request := RegistrationRequest{IsAdd: isAdd, IsOrigin: isOrigin, Datafield: datafield}
	raw, err := proto.Marshal(&request)
	if err != nil {
		return nil, fmt.Errorf("Failed to create RegistrationRequest in NewRegistrationRequest(): %v", err)
	}
	return raw, nil
}

func UnpackRegistrationRequest(raw []byte) (interface{}, error) {
	var request RegistrationRequest
	err := proto.Unmarshal(raw, &request)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal in UnpackRegistrationRequest(): %v", err)
	}
	return &PBRegistrationRequest{request: &request}, nil
}

func NewConnectionRequest(isLogOn bool, swarmID string, originID string, requestCode int) ([]byte, error) {
	request := ConnectionRequest{IsLogOn: isLogOn, SwarmID: swarmID, OriginID: originID, RequestCode: int32(requestCode)}
	raw, err := proto.Marshal(&request)
	if err != nil {
		return nil, fmt.Errorf("Failed to create ConnectionRequest in NewConnectionRequest(): %v", err)
	}
	return raw, nil
}

func UnpackConnectionRequest(raw []byte) (interface{}, error) {
	var request ConnectionRequest
	err := proto.Unmarshal(raw, &request)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal in UnpackConnectionRequest(): %v", err)
	}
	return &PBConnectionRequest{request: &request}, nil
}

func UnmarshalNegotiateMessage(raw []byte) (interface{}, error) {
	var msg NegotiateMessage
	err := proto.Unmarshal(raw, &msg)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal in UnpackConnectionRequest(): %v", err)
	}
	return &PBNegotiateMessage{msg: &msg}, nil
}
