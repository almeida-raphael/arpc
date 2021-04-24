package interfaces

// Serializable Interface is used to identify structs build using Colfer
type Serializable interface {
	MarshalLen()(int, error)
	MarshalTo(buff []byte)int
	UnmarshalBinary(data []byte)error
	Unmarshal(data []byte)(int, error)
}
