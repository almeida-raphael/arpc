package controller

import (
	"context"
	"fmt"
	"github.com/almeida-raphael/arpc/channel"
	"github.com/almeida-raphael/arpc/errors"
	"github.com/almeida-raphael/arpc/headers"
	"github.com/almeida-raphael/arpc/helpers"
	"github.com/almeida-raphael/arpc/interfaces"
	"log"
	"os"
)

// RPC controller struct
type RPC struct {
	channel			  channel.RPC
	logger            *log.Logger

	// Client only
	clientSession     channel.Session

	// Server only
	services		  map[uint32]map[uint16]func(message []byte)([]byte, error)
}

// NewRPCController Creates a new RPC Controller
func NewRPCController(channel channel.RPC) RPC {
	return RPC{
		channel:       channel,
		logger:        log.New(os.Stderr, "aRPC ERROR: ", log.Ldate|log.Ltime),

		clientSession: nil,
		services:      make(map[uint32]map[uint16]func(message []byte)([]byte, error)),
	}
}

func callProcedure(
	stream channel.Stream, header *headers.Header, procedure func(message []byte)([]byte, error),
)([]byte, headers.MessageType, error){
	var messageType headers.MessageType

	data, err := helpers.ReadN(stream, header.PayloadSize)
	if err != nil {
		return nil, 0, err
	}

	result, err := procedure(data)
	if err != nil{
		var errorResponse errors.Error
		errorResponse.Message = fmt.Sprintf("%v", err)

		var errorBytes []byte
		errorResponse.MarshalTo(errorBytes)

		result = errorBytes
		messageType = headers.Error
	}else{
		messageType = headers.Result
	}

	return result, messageType, nil
}

func (c *RPC)_processRemoteCalls(ctx context.Context, session channel.Session){
	stream, err := session.AcceptStream(ctx)
	if err != nil{
		c.logger.Printf( "error accepting stream: %v", err)
		return
	}
	defer func() {
		if err := stream.Close(); err != nil {
			c.logger.Printf( "error closing stream: %v", err)
		}
	}()

	header, err := headers.FromStream(stream)
	if err != nil{
		c.logger.Printf( "%v", err)
		return
	}

	service, ok := c.services[header.ServiceID]
	if ok {
		procedure, ok := service[header.ProcedureID]
		if ok {
			data, messageType, err := callProcedure(stream, header, procedure)
			if err != nil{
				c.logger.Printf( "error calling rpc procedure: %v", err)
				return
			}
			err = c.sendRPCResponse(stream, messageType, header.ServiceID, header.ProcedureID, data)
			if err != nil{
				c.logger.Printf( "error sending rpc response: %v", err)
				return
			}

			return // Success
		}
	}

	c.logger.Printf( "cannot find right procedure to call")
	return
}

func (c *RPC)processRemoteCalls(ctx context.Context, session channel.Session){
	for {
		go c._processRemoteCalls(ctx, session)
	}
}

func (c *RPC) sendRPCResponse(
	stream channel.Stream, messageType headers.MessageType, serviceID uint32, procedureID uint16, request []byte,
)error {
	requestBytes, err := helpers.AddHeaders(messageType, serviceID, procedureID, request)
	if err != nil{
		return err
	}

	err = c.SendData(stream, requestBytes)
	if err != nil{
		return err
	}

	return nil
}

// StartServer starts the RPC server
func (c *RPC)StartServer(ctx context.Context)error{
	listener, err := c.channel.Listen()
	if err != nil{
		return err
	}
	defer func() {
		if err := listener.Close(); err != nil {
			c.logger.Printf( "error closing listener: %v", err)
		}
	}()

	for {
		session, err := listener.Accept(ctx)
		if err != nil {
			c.logger.Printf( "error on accept connection: %v", err)
			continue
		}

		go c.processRemoteCalls(ctx, session)
	}
}

// StartClient starts the RPC client
func (c *RPC)StartClient()error{
	session, err := c.channel.Connect()
	if err != nil{
		return err
	}

	c.clientSession = session
	return nil
}

// RegisterService registers a service to be served on a RPC server
func (c *RPC)RegisterService(serviceID uint32, procedures map[uint16]func(message []byte)([]byte, error)){
	serviceProcedures, ok := c.services[serviceID]
	if !ok{
		serviceProcedures = make(map[uint16]func(message []byte)([]byte, error))
		c.services[serviceID] = serviceProcedures
	}

	for procedureID, procedureFunction := range procedures{
		serviceProcedures[procedureID] = procedureFunction
	}
}

// SendRPCCall sends a RPC message call and deserialize it's response
func (c *RPC)SendRPCCall(
	ctx context.Context, serviceID uint32, procedureID uint16, request interfaces.Serializable,
	response interfaces.Serializable,
)error{
	stream, err := c.clientSession.OpenStream(ctx)
	if err != nil{
		return err
	}
	defer func() {
		if err := stream.Close(); err != nil {
			c.logger.Printf( "error closing stream: %v", err)
		}
	}()

	requestBytes, err := helpers.SerializeWithHeaders(headers.Call, serviceID, procedureID, request)
	if err != nil{
		return err
	}
	err = c.SendData(stream, requestBytes)
	if err != nil{
		return err
	}

	header, err := headers.FromStream(stream)
	if err != nil{
		return err
	}

	responseBytes, err := helpers.ReadN(stream, header.PayloadSize)
	if err != nil{
		return err
	}

	if header.MessageType == uint8(headers.Error){
		var errorResponse errors.Error
		err = errorResponse.UnmarshalBinary(responseBytes)
		if err != nil{
			return err
		}
		return &errors.Remote{Err: fmt.Errorf("%v", errorResponse.Message)}
	}

	err = response.UnmarshalBinary(responseBytes)
	if err != nil{
		return err
	}

	return nil
}

// SendData handles data byte sending through the RPC channel
func (c *RPC)SendData(stream channel.Stream, data []byte)error{
	for len(data) > 0{
		nConsumed, err := stream.Write(data)
		if err != nil{
			return err
		}
		data = data[nConsumed:]
	}

	return nil
}
