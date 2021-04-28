package helpers

import (
	"errors"
	"io"
)

// ReadN reads n bytes from the source reader
func ReadN(reader io.Reader, size uint64)([]byte, error){
	// TODO: Check if this can hang forever
	result := make([]byte, size)
	readBytesCount := 0
	for uint64(len(result)) < size {
		nConsumed, err := reader.Read(result[readBytesCount:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, err
		}
		readBytesCount += nConsumed
	}
	return result, nil
}
