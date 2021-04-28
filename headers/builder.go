package headers

import (
	"fmt"
	"github.com/almeida-raphael/arpc/channel"
	"github.com/almeida-raphael/arpc/helpers"
)

var headerSize = Size()

// BuildHeader generates a new header for a message
func BuildHeader(messageType MessageType, serviceID uint32, procedureID uint16, payloadSize uint64) *Header {
	return &Header{
		MessageType: uint8(messageType),
		ServiceID:   serviceID,
		ProcedureID: procedureID,
		PayloadSize: payloadSize,
	}
}

// FromStream generate a header from bytes from a stream
func FromStream(stream channel.Stream)(*Header, error){
	headerBytes, err := helpers.ReadN(stream, headerSize)
	if err != nil{
		return nil, fmt.Errorf("error on header reading: %v", err)
	}

	var header Header
	err = header.UnmarshalBinary(headerBytes)
	if err != nil{
		return nil, fmt.Errorf("error on header parsing: %v", err)
	}

	return &header, nil
}