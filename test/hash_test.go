package test

import (
	"testing"

	"github.com/mana-sg/vcs/internal/utils"
)

func TestHasher(t *testing.T) {

	testString := []string{
		"hello world!",
		"Test string",
		"test String",
		"test string",
		"Test String",
	}
	for _, str := range testString {
		answer, err := utils.Hash([]byte(str))
		if err != nil {
			t.Errorf("error in result: %v", err)
		}
		t.Logf("String: %s\nHash: %s, Length: %d\n\n", str, string(answer), len(answer))
	}
}
