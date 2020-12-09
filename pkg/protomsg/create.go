package protomsg

import (
	"fmt"

	"github.com/arstevens/go-request/handle"
	"google.golang.org/protobuf/proto"
)

func NewLocalizeRequest(dataspace string, typeCode int) ([]byte, error) {
	request := LocalizeRequest{Type: int32(typeCode), Dataspace: dataspace}
	raw, err := proto.Marshal(&request)
	if err != nil {
		return nil, fmt.Errorf("Failed to create LocalizeRequest in NewLocalizeRequest(): %v", err)
	}
	return raw, nil
}

func UnpackLocalizeRequest(raw []byte, conn handle.Conn) (handle.Request, error) {
	var request LocalizeRequest
	err := proto.Unmarshal(raw, &request)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal in UnpackLocalizeRequest(): %v", err)
	}
	return &PBLocalizeRequest{request: &request}, nil
}

func NewRegistrationRequest(isAdd bool, isOrigin bool, datafield string, typeCode int) ([]byte, error) {
	request := RegistrationRequest{Type: int32(typeCode), IsAdd: isAdd, IsOrigin: isOrigin, Datafield: datafield}
	raw, err := proto.Marshal(&request)
	if err != nil {
		return nil, fmt.Errorf("Failed to create RegistrationRequest in NewRegistrationRequest(): %v", err)
	}
	return raw, nil
}

func UnpackRegistrationRequest(raw []byte, conn handle.Conn) (handle.Request, error) {
	var request RegistrationRequest
	err := proto.Unmarshal(raw, &request)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal in UnpackRegistrationRequest(): %v", err)
	}
	return &PBRegistrationRequest{request: &request}, nil
}

func NewConnectionRequest(isLogOn bool, swarmID string, originID string, requestCode int, typeCode int) ([]byte, error) {
	request := ConnectionRequest{Type: int32(typeCode), IsLogOn: isLogOn, SwarmID: swarmID, OriginID: originID, RequestCode: int32(requestCode)}
	raw, err := proto.Marshal(&request)
	if err != nil {
		return nil, fmt.Errorf("Failed to create ConnectionRequest in NewConnectionRequest(): %v", err)
	}
	return raw, nil
}

func UnpackConnectionRequest(raw []byte, conn handle.Conn) (handle.Request, error) {
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
