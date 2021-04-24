package helpers

import "hash/fnv"

// Hash generates a hash as uint32 to be used on serviceID headers
func Hash(s string) uint32 {
	h := fnv.New32a()

	_, err := h.Write([]byte(s))
	if err != nil { // This error will neve happen
		return 0
	}

	return h.Sum32()
}

