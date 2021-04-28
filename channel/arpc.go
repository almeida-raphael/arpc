package channel

import "context"

// Stream interface definition for a generic aRPC channel
type Stream interface {
	Close()error
	Write([]byte)(int, error)
	Read([]byte)(int, error)
}

// Session interface definition for a generic aRPC channel
type Session interface {
	RemoteAddress()(string, int, error)
	AcceptStream(ctx context.Context)(Stream, error)
	OpenStream(ctx context.Context)(Stream, error)
}

// Listener interface definition for a generic aRPC channel
type Listener interface {
	Close()error
	Accept(context.Context)(Session, error)
}

// RPC interface definition for a generic aRPC channel
type RPC interface {
	Listen()(Listener, error)
	Connect()(Session, error)
}