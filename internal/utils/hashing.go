package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/mana-sg/vcs/pkg/types"
)

func Hash(data []byte) ([]byte, error) {
	if data == nil {
		return nil, fmt.Errorf("nil data cannot be hashed")
	}

	hash := sha256.Sum256(data)

	return []byte(hex.EncodeToString(hash[:])), nil
}

func HashDirectoryContents(children []types.FileNode) ([]byte, error) {
	hasher := sha256.New()

	// Concatenate each child's hash or content
	for _, child := range children {
		var childHash []byte
		if child.Type == 1 {
			childHash, _ = Hash([]byte(child.Content))
		} else {
			childHash, _ = HashDirectoryContents(child.Children)
		}
		hasher.Write(childHash)
	}

	return []byte(hex.EncodeToString(hasher.Sum(nil))), nil
}
