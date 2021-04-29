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

	var responseBytes = make([]byte, uint64(headerSize) + dataSize)
	header.MarshalTo(responseBytes[:headerSize])
	copy(data, responseBytes[headerSize:])

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

	var responseBytes = make([]byte, headerSize + dataSize)
	header.MarshalTo(responseBytes[:headerSize])
	data.MarshalTo(responseBytes[headerSize:])

	return responseBytes, nil
}
