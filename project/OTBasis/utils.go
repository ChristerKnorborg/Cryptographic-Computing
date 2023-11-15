// OTReceiver.go
package project

// Import your ElGamal package
import (
	"math/big"
)

type MessagePair struct {
	Message1 *big.Int
	Message2 *big.Int
}

type PublicKeyPair struct {
	MessageKey0 *big.Int
	MessageKey1 *big.Int
}
