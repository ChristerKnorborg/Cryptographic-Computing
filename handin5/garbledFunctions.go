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

// Create a garbled circuit from the bloodtype compatibility formula (Claudio's Master solution from handin 1)
// For each wire i ∈ [1..T] in the circuit: choose two random strings
func MakeGarbledCircuit() ([][]string, [2]string, [][2]string, [][2]string, [][2]string) {

	// Initialize every wire to two empty strings
	var wires [23][2]string
	for i := 0; i < 22; i++ {
		wires[i] = [2]string{Random128BitString(), Random128BitString()}
	}

	// Create the circuit F from the bloodtype compatibility formula. For every wire, we define a garbled table
	// with the corresponding input and output keys. The first line "XORGate(wires[0], wires[1], wires[2])" represent
	// the left input "wires[0]", the right input "wires[1]" and the output "wires[2]".
	var F [][]string

	// Block 1: x1 and y1
	F = append(F, XORGate(wires[0], wires[1], wires[2])) // XOR constant 1 and x1. Result is ¬x1
	F = append(F, ANDGate(wires[2], wires[3], wires[4])) // AND ¬x1 with y1. Result is z1
	F = append(F, XORGate(wires[4], wires[5], wires[6])) // XOR z1 with constant 1

	// Block 2: x2 and y2
	F = append(F, XORGate(wires[7], wires[8], wires[9]))    // XOR constant 1 and x2. Result is ¬x2
	F = append(F, ANDGate(wires[9], wires[10], wires[11]))  // AND ¬x2 with y2. Result is z2
	F = append(F, XORGate(wires[11], wires[12], wires[13])) // XOR z2 with constant 1

	// Block 3: x3 and y3
	F = append(F, XORGate(wires[14], wires[15], wires[16])) // XOR constant 1 and x3. Result is ¬x3
	F = append(F, ANDGate(wires[16], wires[17], wires[18])) // AND ¬x3 with y3. Result is z3
	F = append(F, XORGate(wires[18], wires[19], wires[20])) // XOR z3 with constant 1

	F = append(F, ANDGate(wires[6], wires[13], wires[21])) // AND z1 and z2. Result is z4

	F = append(F, ANDGate(wires[20], wires[21], wires[22])) // AND z3 and z4. Final

	// Define d = (Z_0, Z_1) = (K^T_0 , K^T_1)
	d := wires[22]

	// Alice's input wires for AND gates
	e_x := [][2]string{wires[1], wires[8], wires[15]}                                   // Alice input bits x1, x2, x3 goes in these wires respectively
	e_y := [][2]string{wires[3], wires[10], wires[17]}                                  // Bob input bits y1, y2, y3 goes in these wires respectively
	e_xor := [][2]string{wires[0], wires[4], wires[7], wires[11], wires[14], wires[18]} // constants from XOR gates

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

// The garbled evaluation function Ev that evaluates a garbled circuit F on a garbled input X and produces a garbled output Z′
func GarbledEvaluation() {

}

// Encoding function En uses e to map Bob's input y to a garbled input Y, by encoding each bit of y into the corresponding label
func Encode(e_y [][2]string, y int) []string {

	// Make an array of the three bits of Bob's input y
	inputInBits := extractBits(y)

	var Y []string

	for i := 0; i < 3; i++ {
		// Match the input bit with the corresponding label
		Y = append(Y, e_y[i][inputInBits[i]])
	}

	return Y
}

// The decoding function De, decodes the garbled output Z′ into a plaintext output z.
func Decode() {

}

// This function extract the three bits from an bloodtype input and returns them as an array
func extractBits(n int) [3]int {
	// Extract the bits from Bob (donor)

	n1 := (n >> 2) & 1 // extract 3rd rightmost bit
	n2 := (n >> 1) & 1 // extract 2nd rightmost bit
	n3 := n & 1        // extract rightmost bit

	return [3]int{n1, n2, n3}
}
