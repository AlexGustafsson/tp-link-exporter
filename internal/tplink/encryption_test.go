package tplink

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncryptionRoundtrip(t *testing.T) {
	message := []byte("hello, world")
	decrypted := DecryptMessage(EncryptMessage(message))
	assert.Equal(t, message, decrypted)
}
