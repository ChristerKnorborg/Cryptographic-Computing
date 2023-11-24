package utils

// Import your ElGamal package
import (
	"crypto/rand"
	"crypto/sha256"
	"cryptographic-computing/project/elgamal"
	"fmt"
	"log"
	"math/big"
	mathRand "math/rand"
	"strconv"
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

type Seed struct {
	Seed0 *big.Int
	Seed1 *big.Int
}

type ByteCiphertextPair struct {
	Y0 []byte
	Y1 []byte
}

func PseudoRandomGenerator(seed *big.Int, bitLength int) ([]uint8, error) {
	if bitLength <= 0 {
		return nil, fmt.Errorf("bitLength must be positive")
	}
	output := make([]uint8, 0, bitLength) // Allocate space for the array

	// Convert seed to a byte slice
	seedBytes := seed.Bytes()

	for len(output) < bitLength {
		hash := sha256.Sum256(seedBytes)
		for _, b := range hash[:] {
			for i := 0; i < 8 && len(output) < bitLength; i++ {
				bit := (b >> (7 - i)) & 1 // Extract each bit
				output = append(output, uint8(bit))
			}
		}

		// Increment the seed
		seed = new(big.Int).Add(seed, big.NewInt(1))
		seedBytes = seed.Bytes()
	}

	return output, nil
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
		case uint8:
			xorResult ^= int(v)
		default:
			return "", fmt.Errorf("unsupported type")
		}
	}

	return strconv.Itoa(xorResult), nil
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
