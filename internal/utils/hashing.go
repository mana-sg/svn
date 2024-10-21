package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func Hash(data []byte) ([]byte, error) {
	if data == nil {
		return nil, fmt.Errorf("nil data cannot be hashed")
	}

	hash := sha256.Sum256(data)

	return []byte(hex.EncodeToString(hash[:])), nil
}
