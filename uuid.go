package indexstore

import "github.com/google/uuid"

func uuidFromBytes(b []byte) uuid.UUID {
	if len(b) != 16 {
		return uuid.Nil
	}

	var dst uuid.UUID
	copy(dst[:], b)
	return dst
}
