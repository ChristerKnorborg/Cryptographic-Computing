// OTReceiver.go
package OTExtension

// Import your ElGamal package
import (
	"crypto/rand"
	"crypto/sha256"
	"cryptographic-computing/project/elgamal"
	"encoding/hex"
	"math/big"
)

// Struct to store messages M0 and M1
type MessagePair struct {
	Message0 string
	Message1 string
}

// Struct to store public keys for Oblivious transfer. Each public key pair consists of 2 public keys -
// one real key and one oblivious key.
type PublicKeyPair struct {
	MessageKey0 *big.Int
	MessageKey1 *big.Int
}

// Struct to store ciphertexts the two messages M0 and M1
type CiphertextPair struct {
	Ciphertext0 *elgamal.Ciphertext
	Ciphertext1 *elgamal.Ciphertext
}

type seed struct {
	seed0 *big.Int
	seed1 *big.Int
}

func pseudoRandomGenerator(seed *big.Int, length int) ([]byte, error) {
	// For the sake of example, we're generating a random bit string.
	// Replace this with an actual PRG function in production.
	output := make([]byte, length)
	_, err := rand.Read(output)
	if err != nil {
		return nil, err
	}
	return output, nil
}

// HashFunction is an example of how to hash data using SHA-256.
func HashFunction(data []byte) string {
	hasher := sha256.New()
	hasher.Write(data)
	return hex.EncodeToString(hasher.Sum(nil))
}
