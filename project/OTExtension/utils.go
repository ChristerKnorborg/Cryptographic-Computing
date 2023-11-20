// OTReceiver.go
package OTExtension

// Import your ElGamal package
import (
	"crypto/sha256"
	"cryptographic-computing/project/elgamal"
	"encoding/hex"
	"math/big"
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

type seed struct {
	seed0 *big.Int
	seed1 *big.Int
}

func pseudoRandomGenerator(seed *big.Int, bitLength int) ([]byte, error) {
	// Convert the length in bits to length in bytes, rounding up
	byteLength := (bitLength + 7) / 8
	output := make([]byte, 0, byteLength)

	// Convert seed to a byte slice
	seedBytes := seed.Bytes()

	// Hash the seed and append to output until we have enough bytes
	for len(output) < byteLength {
		hash := sha256.Sum256(seedBytes)
		output = append(output, hash[:]...)

		// Increment the seed to simulate feeding the output of the PRG back into it
		seed = new(big.Int).Add(seed, big.NewInt(1))
		seedBytes = seed.Bytes()
	}

	// Trim the output to the exact number of bytes we need
	output = output[:byteLength]

	// If the bitLength is less than 8, we need to extract the exact number of bits.
	if bitLength < 8 {
		// Take the first 8 bits from output[0] and shift right to get only the bits we need.
		return []byte{output[0] >> (8 - bitLength)}, nil
	}

	// If the bitLength is not a multiple of 8, we need to clear the excess bits in the last byte
	if excessBits := 8*byteLength - bitLength; excessBits > 0 {
		// Clear the excess bits
		output[byteLength-1] &= 0xFF << excessBits
	}

	return output, nil
}

// HashFunction is an example of how to hash data using SHA-256.
func HashFunction(data []byte) string {
	hasher := sha256.New()
	hasher.Write(data)
	return hex.EncodeToString(hasher.Sum(nil))
}

func PrintMatrix(matrix [][]byte) {
	for i := 0; i < len(matrix); i++ {
		for j := 0; j < len(matrix[0]); j++ {

			print(matrix[i][j], " ")
		}
		println()
	}
	println()
}
