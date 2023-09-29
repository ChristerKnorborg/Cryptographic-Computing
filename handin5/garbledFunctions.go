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

// Make the garbled circuit of the bloodtype compatibility formula
func MakeGarbledCircuit(y int) ([][]string, [2]string, [][2]string, [][2]string, [][2]string) {

	// Create two labels for each wire in the circuit
	var labels [23][2]string
	for i := 0; i < 22; i++ {
		labels[i] = [2]string{"", ""} // Initialize the labels to empty
	}

	var F [][]string

	// Block 1: x1 and y1
	F = append(F, XORGate(labels[0], labels[1], labels[2])) // XOR constant 1 and x1. Result is ¬x1
	F = append(F, ANDGate(labels[2], labels[3], labels[4])) // AND ¬x1 with y1. Result is z1
	F = append(F, XORGate(labels[4], labels[5], labels[6])) // XOR z1 with constant 1

	// Block 2: x2 and y2
	F = append(F, XORGate(labels[7], labels[8], labels[9]))    // XOR constant 1 and x2. Result is ¬x2
	F = append(F, ANDGate(labels[9], labels[10], labels[11]))  // AND ¬x2 with y2. Result is z2
	F = append(F, XORGate(labels[11], labels[12], labels[13])) // XOR z2 with constant 1

	// Block 3: x3 and y3
	F = append(F, XORGate(labels[14], labels[15], labels[16])) // XOR constant 1 and x3. Result is ¬x3
	F = append(F, ANDGate(labels[16], labels[17], labels[18])) // AND ¬x3 with y3. Result is z3
	F = append(F, XORGate(labels[18], labels[19], labels[20])) // XOR z3 with constant 1

	F = append(F, ANDGate(labels[6], labels[13], labels[21])) // AND z1 and z2. Result is z4

	F = append(F, ANDGate(labels[20], labels[21], labels[22])) // AND z3 and z4. Final

	d := labels[22]

	// Alice's input labels for AND gates
	e_x := [][2]string{labels[1], labels[8], labels[15]}                                      // Alice input bits x1, x2, x3 goes in these labels respectively
	e_y := [][2]string{labels[3], labels[10], labels[17]}                                     // Bob input bits y1, y2, y3 goes in these labels respectively
	e_xor := [][2]string{labels[0], labels[4], labels[7], labels[11], labels[14], labels[18]} // constants from XOR gates

	return F, d, e_x, e_y, e_xor
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
