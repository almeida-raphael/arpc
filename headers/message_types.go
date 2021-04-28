package headers

// MessageType possible types for the message type on the header
type MessageType uint8

// Available message types
const (
	Result MessageType = iota
	Call
	Error
)