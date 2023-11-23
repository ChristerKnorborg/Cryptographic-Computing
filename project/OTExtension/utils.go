// OTReceiver.go
package OTExtension

// Import your ElGamal package
import (
	"crypto/sha256"
	"cryptographic-computing/project/elgamal"
	"fmt"
	"math/big"
	"strconv"
	"strings"
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

type Seed struct {
	seed0 *big.Int
	seed1 *big.Int
}

type ByteCiphertextPair struct {
	y0 []byte
	y1 []byte
}

func pseudoRandomGenerator(seed *big.Int, bitLength int) (string, error) {
	// Convert the length in bits to length in bytes, rounding up
	byteLength := (bitLength + 7) / 8
	output := make([]byte, 0, byteLength)

	// Convert seed to a byte slice
	seedBytes := seed.Bytes()

	// Hash the seed and append to output until we have enough bytes
	for len(output) < byteLength {
		hash := sha256.Sum256(seedBytes)
		output = append(output, hash[:]...)

		// Increment the seed
		seed = new(big.Int).Add(seed, big.NewInt(1))
		seedBytes = seed.Bytes()
	}

	// Trim the output to the exact number of bytes we need
	output = output[:byteLength]

	// Convert the bytes to a binary string
	var stringBuilder strings.Builder
	for _, b := range output {
		stringBuilder.WriteString(fmt.Sprintf("%08b", b))
	}

	// If the bitLength is not a multiple of 8, trim the excess bits from the end of the string
	bitString := stringBuilder.String()
	if excessBits := 8*byteLength - bitLength; excessBits > 0 {
		bitString = bitString[:len(bitString)-excessBits]
	}

	return bitString, nil
}

// Hash creates a hash of the input data with a specified byte length.
func Hash(data []byte, byteLength int) []byte {

	fullHash := make([]byte, 0, byteLength)

	// SHA-256 produces a hash of 32 bytes (256 bits)
	// We keep hashing and concatenating until we reach the desired byte length
	for len(fullHash) < byteLength {
		hash := sha256.Sum256(data)
		fullHash = append(fullHash, hash[:]...)

		// Modify data slightly for the next iteration to produce a different hash
		// For example, append a byte that represents the current length of fullHash
		data = append(data, byte(len(fullHash)))
	}

	// Truncate the hash to the exact byte length required
	return fullHash[:byteLength]
}

func PrintMatrix(matrix [][]string) {
	for i := 0; i < len(matrix); i++ {
		for j := 0; j < len(matrix[0]); j++ {

			print(matrix[i][j], " ")
		}
		println()
	}
	println()
}

// PrintBinaryString prints the binary representation of a []byte slice.
func PrintBinaryString(bytes []byte) {
	binaryString := ""
	for _, b := range bytes {
		binaryString += fmt.Sprintf("%08b", b) // Convert each byte to an 8-bit binary string
	}
	fmt.Println("As binary string:", binaryString)
}

// XOR takes a variable number of arguments (strings and ints),
// performs a bitwise XOR operation on all of them, and returns the result as a string.
func XOR(args ...interface{}) (string, error) {
	var xorResult int

	for _, arg := range args {
		switch v := arg.(type) {
		case string:
			val, err := strconv.Atoi(v)
			if err != nil {
				return "", err
			}
			xorResult ^= val
		case int:
			xorResult ^= v
		default:
			return "", fmt.Errorf("unsupported type")
		}
	}

	return strconv.Itoa(xorResult), nil
}
