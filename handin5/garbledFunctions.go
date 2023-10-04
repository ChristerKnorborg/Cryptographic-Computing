package handin5

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
)

// Struct used for the key values for a single wire in the garbled circuit
type KeyPair struct {
	K_0 string //
	K_1 string
}

// Struct used for the garbled gate with 4 truth table entries randomly permuted
type GarbledGate struct {
	C_0 string // Truth table entry 0
	C_1 string // Truth table entry 1
	C_2 string // Truth table entry 2
	C_3 string // Truth table entry 3
}

// Create a garbled gate for the XOR gate
func XORGate(inputL KeyPair, inputR KeyPair, output KeyPair) GarbledGate {

	stringOfZeros := Zeros128BitString() // Create a string of 128 zeros
	gg := GarbledGate{}                  // Initialize the garbled gate

	// Hash 256 bits of two input 128-bit keys concatenated
	c1_hash := Hash(inputL.K_0, inputR.K_0)
	c2_hash := Hash(inputL.K_0, inputR.K_1)
	c3_hash := Hash(inputL.K_1, inputR.K_0)
	c4_hash := Hash(inputL.K_1, inputR.K_1)

	// Add redundancy (later used to check if decryption is successful)
	c1_concat := output.K_0 + stringOfZeros // output is 0, as XOR(0,0) = 0
	c2_concat := output.K_1 + stringOfZeros // output is 1, as XOR(0,1) = 1
	c3_concat := output.K_1 + stringOfZeros // output is 1, as XOR(1,0) = 1
	c4_concat := output.K_0 + stringOfZeros // output is 0, as XOR(1,1) = 0

	// Encrypt the output key(s) with the hash of the input keys
	gg.C_0 = XORStrings(c1_hash, c1_concat)
	gg.C_1 = XORStrings(c2_hash, c2_concat)
	gg.C_2 = XORStrings(c3_hash, c3_concat)
	gg.C_3 = XORStrings(c4_hash, c4_concat)

	// Randomly shuffle the encrypted output keys
	gg.Shuffle()
	return gg
}

// Create a garbled gate for the AND gate
func ANDGate(inputL KeyPair, inputR KeyPair, output KeyPair) GarbledGate {

	stringOfZeros := Zeros128BitString() // Create a string of 128 zeros
	gg := GarbledGate{}                  // Initialize the garbled gate

	// Hash 256 bits of two input 128-bit keys concatenated
	c1_hash := Hash(inputL.K_0, inputR.K_0)
	c2_hash := Hash(inputL.K_0, inputR.K_1)
	c3_hash := Hash(inputL.K_1, inputR.K_0)
	c4_hash := Hash(inputL.K_1, inputR.K_1)

	// Add redundancy (later used to check if decryption is successful)
	c1_concat := output.K_0 + stringOfZeros // output is 0, as AND(0,0) = 0
	c2_concat := output.K_0 + stringOfZeros // output is 0, as AND(0,1) = 0
	c3_concat := output.K_0 + stringOfZeros // output is 0, as AND(1,0) = 0
	c4_concat := output.K_1 + stringOfZeros // output is 1, as AND(1,1) = 1

	// Encrypt the output key(s) with the hash of the input keys
	gg.C_0 = XORStrings(c1_hash, c1_concat)
	gg.C_1 = XORStrings(c2_hash, c2_concat)
	gg.C_2 = XORStrings(c3_hash, c3_concat)
	gg.C_3 = XORStrings(c4_hash, c4_concat)

	// Randomly shuffle the encrypted output keys
	gg.Shuffle()

	return gg
}

// Hash concatenates two strings and hashes them using SHA256. The output is a 256-bit hash as a 64 char hex string
func Hash(leftKey string, rightKey string) string {
	hasher := sha256.New()
	hasher.Write([]byte(leftKey))
	hasher.Write([]byte(rightKey))

	return hex.EncodeToString(hasher.Sum(nil))
}

// Shuffle permutes the entries of a given GarbleGate struct randomly
func (gg *GarbledGate) Shuffle() {

	// Create a slice of the truth table entries
	entries := []*string{&gg.C_0, &gg.C_1, &gg.C_2, &gg.C_3}

	// Use the Fisher-Yates shuffle
	for i := len(entries) - 1; i > 0; i-- {
		j, _ := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		*entries[i], *entries[j.Int64()] = *entries[j.Int64()], *entries[i]
	}
	// Reassign the shuffled values back to the struct
	gg.C_0, gg.C_1, gg.C_2, gg.C_3 = *entries[0], *entries[1], *entries[2], *entries[3]
}

// Creates a string of 128 zeros (hex-encoded as 32 characters of 0)
func Zeros128BitString() string {
	return strings.Repeat("0", 32) // 32 hex characters = 128 bits
}

// Creates a random 128-bit string (hex-encoded as 32 characters)
func Random128BitString() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}
	return hex.EncodeToString(bytes)

}

// XORStrings XORs two hex-encoded strings and returns the result as a hex-encoded string.
func XORStrings(a string, b string) string {
	bytesA, errA := hex.DecodeString(a)
	bytesB, errB := hex.DecodeString(b)
	if errA != nil || errB != nil || len(bytesA) != len(bytesB) {
		fmt.Println("a: ", a)
		fmt.Println("b: ", b)
		panic("invalid input or length mismatch")
	}

	result := make([]byte, len(bytesA))
	for i := range bytesA {
		result[i] = bytesA[i] ^ bytesB[i]
	}

	return hex.EncodeToString(result)
}

// This function extract the three bits from an bloodtype input and returns them as an array
func ExtractBits(n int) [3]int {

	n1 := (n >> 2) & 1 // extract 3rd rightmost bit
	n2 := (n >> 1) & 1 // extract 2nd rightmost bit
	n3 := n & 1        // extract rightmost bit

	return [3]int{n1, n2, n3}
}

// EvaluateGarbledGate evaluates a single garbled gate using the keys from the input wires.
// It returns the output key of the garbled gate as a string.
func EvaluateGarbledGate(gate GarbledGate, leftKey string, rightKey string) string {

	zeroString := Zeros128BitString()
	hashValue := Hash(leftKey, rightKey)

	// List of garbled gate entries
	gateEntries := []string{gate.C_0, gate.C_1, gate.C_2, gate.C_3}

	debugString := []string{} // Used for debugging

	for i, entry := range gateEntries {

		decryptedTableEntry := XORStrings(entry, hashValue)

		debugString = append(debugString, fmt.Sprintf("Decrypted table entry %d: %s\n", i, decryptedTableEntry)) // Used for debugging

		// Check if the last 128 bits of the decrypted table entry is a string of 128 zeros.
		// If it is, return the first 128 bits of the decrypted table entry.
		if strings.HasSuffix(decryptedTableEntry, zeroString) {
			return decryptedTableEntry[:32]
		}
	}
	fmt.Println(debugString) // Used for debugging
	panic("Decryption failed. The decrypted table entry does not end with a string of 128 zeros. The 4 results are above.")
}
