// OTReceiver.go
package OTBasic

// Import your ElGamal package
import (
	"crypto/rand"
	"cryptographic-computing/project/elgamal"
	"log"
	"math/big"
	mathRand "math/rand"
	"time"
)

// Struct to store messages M0 and M1
type MessagePair struct {
	Message0 []byte
	Message1 []byte
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

// RandomBytes generates a slice of random bytes of a given length
func RandomBytes(length int) []byte {
	b := make([]byte, length)
	_, err := rand.Read(b)
	// Handle the error here. In production code, you might want to pass it up the call stack
	if err != nil {
		log.Fatal(err)
	}
	return b
}

// AllOnesBytes generates a slice of bytes of a given length, all set to 1
func AllOnesBytes(length int) []byte {
	b := make([]byte, length)
	for i := range b {
		b[i] = 1
	}
	return b
}

// AllTwosBytes generates a slice of bytes of a given length, all set to 2
func AllTwosBytes(length int) []byte {
	b := make([]byte, length)
	for i := range b {
		b[i] = 2
	}
	return b
}

// RandomSelectionBits generates a slice of m random selection bits (0 or 1)
func RandomSelectionBits(m int) []int {

	mathRand.NewSource(time.Now().UnixNano())

	bits := make([]int, m)
	for i := range bits {
		bits[i] = mathRand.Intn(2) // Generates a random integer 0 or 1
	}
	return bits
}
