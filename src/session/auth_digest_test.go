package session

import (
	"crypto/sha1"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAuthDigest(t *testing.T) {
	testCases := map[string]string{
		"Notch": "4ed1f46bbe04bc756bcb17c0c7ce3e4632f06a48",
		"jeb_":  "-7c9d5b0044c130109a5d7b5fb5c317c02b4e28c1",
		"simon": "88e16a1019277b15d58faf0541e11910eb756f6",
	}

	for s, hash := range testCases {
		t.Run(s, func(t *testing.T) {
			h := sha1.New()
			_, _ = h.Write([]byte(s))
			sum := h.Sum(nil)

			digest := AuthDigest(sum)
			assert.Equal(t, hash, digest)
		})
	}
}
