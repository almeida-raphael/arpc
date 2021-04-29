package headers

import (
	"github.com/almeida-raphael/arpc/interfaces"
)

// AddHeaders add headers to already serialized data
func AddHeaders(messageType MessageType, serviceID uint32, procedureID uint16, data []byte)([]byte, error){
	dataSize := uint64(len(data))

	header := BuildHeader(messageType, serviceID, procedureID, dataSize)
	headerSize, err := header.MarshalLen()
	if err != nil {
		return nil, err
	}

	var responseBytes = make([]byte, 1 + uint64(headerSize) + dataSize)
	responseBytes[0] = uint8(headerSize)
	header.MarshalTo(responseBytes[1:headerSize+1])
	copy(data, responseBytes[headerSize+1:])

	return responseBytes, nil
}


// SerializeWithHeaders Serialize a given message with it's headers
func SerializeWithHeaders(
	messageType MessageType, serviceID uint32, procedureID uint16, data interfaces.Serializable,
)([]byte, error){
	dataSize, err := data.MarshalLen()
	if err != nil {
		return nil, err
	}

	header := BuildHeader(messageType, serviceID, procedureID, uint64(dataSize))
	headerSize, err := header.MarshalLen()
	if err != nil {
		return nil, err
	}

	var responseBytes = make([]byte, 1 + headerSize + dataSize)
	responseBytes[0] = uint8(headerSize)
	header.MarshalTo(responseBytes[1:headerSize+1])
	data.MarshalTo(responseBytes[headerSize+1:])

	return responseBytes, nil
}
