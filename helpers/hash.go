package helpers

import "hash/fnv"

func Hash(s string) uint32 {
	h := fnv.New32a()

	_, err := h.Write([]byte(s))
	if err != nil { // This error will neve happen
		return 0
	}

	return h.Sum32()
}

