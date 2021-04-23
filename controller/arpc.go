package controller

import (
	"github.com/almeida-raphael/aRPC/helpers"
	"github.com/almeida-raphael/aRPC/interfaces"
)

type RPC struct {}

func (rpcController *RPC)RegisterService(serviceID uint32, procedures map[uint16]func(message []byte)([]byte, error))error{
	return nil
}

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

	// TODO: Consume Response

	return nil
}

func (rpcController *RPC)SendData(data []byte)error{
	return nil
}