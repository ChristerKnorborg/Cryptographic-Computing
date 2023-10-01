package handin5

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
)

// Create a garbled table for the XOR gate
func XORGate(inputL [2]string, inputR [2]string, output [2]string) []string {

	stringOfZeros := Zeros128BitString() // Create a string of 128 zeros
	gg := make([]string, 4)              // Initialize 4 truth table entries

	// Hash 256 bits of input 128-bit keys concatenated
	c1_hash := Hash(inputL[0], inputR[0])
	c2_hash := Hash(inputL[0], inputR[1])
	c3_hash := Hash(inputL[1], inputR[0])
	c4_hash := Hash(inputL[1], inputR[1])

	// Add redundancy (later used to check if decryption is successful)
	c1_concat := output[0] + stringOfZeros // output is 0, as XOR(0,0) = 0
	c2_concat := output[1] + stringOfZeros // output is 1, as XOR(0,1) = 1
	c3_concat := output[1] + stringOfZeros // output is 1, as XOR(1,0) = 1
	c4_concat := output[0] + stringOfZeros // output is 0, as XOR(1,1) = 0

	// Encrypt the output key with the input keys
	gg = append(gg, XORStrings(c1_hash, c1_concat))
	gg = append(gg, XORStrings(c2_hash, c2_concat))
	gg = append(gg, XORStrings(c3_hash, c3_concat))
	gg = append(gg, XORStrings(c4_hash, c4_concat))

	// Randomly shuffle the ciphertexts (to hide information about inputs/outputs)
	gg = Shuffle(gg)

	return gg
}

// Create a garbled table for the AND gate
func ANDGate(inputL [2]string, inputR [2]string, output [2]string) []string {

	stringOfZeros := Zeros128BitString() // Create a string of 128 zeros
	gg := make([]string, 4)              // Initialize 4 truth table entries

	// Hash 256 bits of input 128-bit keys concatenated
	c1_hash := Hash(inputL[0], inputR[0])
	c2_hash := Hash(inputL[0], inputR[1])
	c3_hash := Hash(inputL[1], inputR[0])
	c4_hash := Hash(inputL[1], inputR[1])

	// Add redundancy (later used to check if decryption is successful)
	c1_concat := output[0] + stringOfZeros // output is 0, as AND(0,0) = 0
	c2_concat := output[0] + stringOfZeros // output is 0, as AND(0,1) = 0
	c3_concat := output[0] + stringOfZeros // output is 0, as AND(1,0) = 0
	c4_concat := output[1] + stringOfZeros // output is 1, as AND(1,1) = 1

	// Encrypt the output key with the input keys
	gg = append(gg, XORStrings(c1_hash, c1_concat))
	gg = append(gg, XORStrings(c2_hash, c2_concat))
	gg = append(gg, XORStrings(c3_hash, c3_concat))
	gg = append(gg, XORStrings(c4_hash, c4_concat))

	// Randomly shuffle the ciphertexts (to hide information about inputs/outputs)
	gg = Shuffle(gg)

	return gg
}

// Hash concatenates two strings and hashes them using SHA256. The output is a 256-bit hash as a hex string.
func Hash(leftKey string, rightKey string) string {
	hasher := sha256.New()
	hasher.Write([]byte(leftKey))
	hasher.Write([]byte(rightKey))

	return fmt.Sprintf("%x", hasher.Sum(nil))
}

// Shuffle shuffles the given slice of strings using the Fisher-Yates algorithm.
// This algorithm is taken directly from ChatGPT, since we found out the standard library
// shuffle function is not cryptographically secure.
func Shuffle(slice []string) []string {
	for i := len(slice) - 1; i > 0; i-- {
		j, _ := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		slice[i], slice[j.Int64()] = slice[j.Int64()], slice[i]
	}
	return slice
}

// Creates a string of 128 zeros
func Zeros128BitString() string {
	return strings.Repeat("0", 128)
}

// Creates a random 128-bit string
func Random128BitString() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%x", b)
}

// XORStrings XORs two hex-encoded strings and returns the result as a hex-encoded string.
func XORStrings(a string, b string) string {
	bytesA, errA := hex.DecodeString(a)
	bytesB, errB := hex.DecodeString(b)
	if errA != nil || errB != nil || len(bytesA) != len(bytesB) {
		panic("invalid input or length mismatch")
	}

	result := make([]byte, len(bytesA))
	for i := range bytesA {
		result[i] = bytesA[i] ^ bytesB[i]
	}

	return hex.EncodeToString(result)
}

func CalculateGate(gate []string, leftKey string, rightKey string) string {
	// Hash 256 bits of input 128-bit keys concatenated
	c1_hash := Hash(leftKey, rightKey)

	// Decrypt the output key with the input keys
	return XORStrings(c1_hash, gate[0])
}

// The garbled evaluation function Ev that evaluates a garbled circuit F on a garbled input X and produces a garbled output Z′
func GarbledEvaluation() {

}

// The decoding function De, decodes the garbled output Z′ into a plaintext output z.
func Decode() {

}

// This function extract the three bits from an bloodtype input and returns them as an array
func ExtractBits(n int) [3]int {
	// Extract the bits from Bob (donor)

	n1 := (n >> 2) & 1 // extract 3rd rightmost bit
	n2 := (n >> 1) & 1 // extract 2nd rightmost bit
	n3 := n & 1        // extract rightmost bit

	return [3]int{n1, n2, n3}
}

// EvaluateGarbledGate evaluates a single garbled gate using the keys from the OT
func EvaluateGarbledGate(gate []string, leftKey string, rightKey string) string {

	zeroStrng := Zeros128BitString()

	// Try all truth table entries that are randomly shuffled
	for i := 0; i < 4; i++ {
		decryptedTableEntry := XORStrings(gate[i], Hash(leftKey, rightKey))

		// Check if the last 128 bits of the decrypted table entry is a string of 128 zeros
		if strings.HasSuffix(decryptedTableEntry, zeroStrng) {

			// return the first 128 bits of the decrypted table entry.
			// Notice, that this get the correct truth table entry with overwhelming probability (but not with certainty)
			return decryptedTableEntry[:128]
		}

	}

}
