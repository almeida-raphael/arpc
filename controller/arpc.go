package controller

import (
	"github.com/almeida-raphael/arpc/helpers"
	"github.com/almeida-raphael/arpc/interfaces"
)

// RPC controller struct
type RPC struct {}

// RegisterService registers a service to be served on a RPC server
func (rpcController *RPC)RegisterService(serviceID uint32, procedures map[uint16]func(message []byte)([]byte, error))error{
	return nil
}

// SendRPC sends a RPC message and deserialize it's response
func (rpcController *RPC)SendRPC(
	messageType uint8, serviceID uint32, procedureID uint16, request interfaces.Serializable,
	response interfaces.Serializable,
)error{
	requestBytes, err := helpers.SerializeWithHeaders(messageType, serviceID, procedureID, request)
	if err != nil{
		return err
	}

	err = rpcController.SendData(requestBytes)
	if err != nil{
		return err
	}

	// TODO: Consume Response if message type is != Result

	return nil
}

// SendData handles data byte sending through the RPC channel
func (rpcController *RPC)SendData(data []byte)error{
	return nil
}