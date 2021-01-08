package debriefer

import (
	"encoding/binary"
	"io"
	"log"
)

func LoadPreferenceDebrief(conn io.Reader) interface{} {
	var debriefValue int32
	err := binary.Read(conn, binary.BigEndian, &debriefValue)
	if err != nil {
		log.Printf("Failed to debrief in gateway.debriefConnection(): %v", err)
		return -1
	}
	return int(debriefValue)
}
