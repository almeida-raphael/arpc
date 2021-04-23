package helpers

import (
	"github.com/almeida-raphael/aRPC/headers"
	"github.com/almeida-raphael/aRPC/interfaces"
)

func SerializeWithHeaders(
	messageType uint8, serviceID uint32, procedureID uint16, data interfaces.Serializable,
)([]byte, error){
	header, err := headers.BuildHeader(messageType, serviceID, procedureID)
	if err != nil {
		return nil, err
	}

	headerSize, err := header.MarshalLen()
	if err != nil {
		return nil, err
	}

	dataSize, err := data.MarshalLen()
	if err != nil {
		return nil, err
	}

	var responseBytes = make([]byte, headerSize + dataSize)

	header.MarshalTo(responseBytes[0:headerSize])
	data.MarshalTo(responseBytes[headerSize:dataSize])

	return responseBytes, nil
}

func DeserializeHeader(data []byte)(*headers.Header, []byte, error){
	var header headers.Header

	consumedHeader, err := header.Unmarshal(data)
	if err != nil {
		return nil, data, err
	}

	return &header, data[consumedHeader:], nil
}