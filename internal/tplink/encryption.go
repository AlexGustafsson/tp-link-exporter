package tplink

// EncryptMessage encrypts a message. Each byte is XOR'd with the previous encrypted byte.
// The first byte is XOR'd with 0xab.
func EncryptMessage(message []byte) []byte {
	key := uint8(0xab)
	encrypted := make([]byte, len(message))
	for i := 0; i < len(message); i++ {
		encrypted[i] = byte(uint8(message[i]) ^ key)
		key = encrypted[i]
	}
	return encrypted
}

// DecryptMessage decrypts a message. Each byte is XOR'd with the previous encrypted byte.
// The first byte is XOR'd with 0xab.
func DecryptMessage(message []byte) []byte {
	key := uint8(0xab)
	decrypted := make([]byte, len(message))
	for i := 0; i < len(message); i++ {
		nextKey := message[i]
		decrypted[i] = byte(uint8(message[i]) ^ key)
		key = nextKey
	}
	return decrypted
}
