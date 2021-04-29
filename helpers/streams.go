package helpers

import (
	"errors"
	"io"
)

// ReadN reads n bytes from the source reader
func ReadN(reader io.Reader, size uint64)([]byte, error){
	// TODO: Check if this can hang forever
	result := make([]byte, size)
	var readBytesCount uint64 = 0
	for readBytesCount < size {
		nConsumed, err := reader.Read(result[readBytesCount:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, err
		}
		readBytesCount += uint64(nConsumed)
	}
	return result, nil
}
