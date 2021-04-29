package helpers

import "hash/fnv"

// Hash generates a hash as uint32 to be used on serviceID headers
func Hash(s string) uint32 {
	h := fnv.New32a() // TODO: Check if this is unique on every machine

	_, err := h.Write([]byte(s))
	if err != nil { // This error will never happen
		return 0
	}

	return h.Sum32()
}

