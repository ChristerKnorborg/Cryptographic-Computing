package project

import (
	"crypto/cipher"
)

type OTSender struct {
	k0, k1 [32]byte // keys for two messages of 256 bits
}

// OTRReceiver is the receiver in the OT protocol.
type OTRReceiver struct {
	sigma byte         // choice bit: 0 or 1
	k     [32]byte     // chosen key
	c     cipher.Block // chosen cipher
}

// NewOTSender creates a new sender with two messages.
func NewOTSender(m0, m1 string, elGamal *ElGamal) (*OTSender, [2][]byte, error) {
	// Generate random keys
	var k0, k1 [32]byte

	// Generate random keys
	elGamal

	// Encrypt messages with elGamal

}
