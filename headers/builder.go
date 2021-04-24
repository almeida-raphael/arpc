package headers

import (
	"encoding/binary"
	"github.com/google/uuid"
)

// BuildHeader generates a new header for a message
func BuildHeader(messageType uint8, serviceID uint32, procedureID uint16)(*Header, error){
	uuidBinary, err := uuid.New().MarshalBinary()
	if err != nil {
		return nil, err
	}

	header := Header{
		ID:          binary.BigEndian.Uint16(uuidBinary),
		MessageType: messageType,
		ServiceID:   serviceID,
		ProcedureID: procedureID,
	}

	return &header, nil
}
