package handin5

import "math/big"

type Bob struct {
	y          int        // Bob input
	publicKeys []*big.Int // Public keys from Alice
	F          [][]string // Garbled circuit
	d          [2]string  // The Z values from the ouput of the garbled circuit
	e_x        [][2]string
	e_y        [][2]string
	e_xor      [][2]string
}

type Y struct {
	encoded_y []string
	e_xor     [][2]string
}

// Set Bob's input as the y provided by the GarbledCircuit function
func (bob *Bob) Init(y int) {
	bob.y = y
}

func (bob *Bob) ReceiveKeys(publicKeys []*big.Int) {
	bob.publicKeys = publicKeys
}

// Create a garbled circuit from the bloodtype compatibility formula (Claudio's Master solution from handin 1)
// For each wire i ∈ [1..T] in the circuit: choose two random strings
func (bob *Bob) MakeGarbledCircuit() ([][]string, [2]string, [][2]string) {

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

	e_x := [][2]string{wires[1], wires[8], wires[15]}                                   // Alice input bits x1, x2, x3 goes in these wires respectively
	e_y := [][2]string{wires[3], wires[10], wires[17]}                                  // Bob input bits y1, y2, y3 goes in these wires respectively
	e_xor := [][2]string{wires[0], wires[4], wires[7], wires[11], wires[14], wires[18]} // constants from XOR gates goes in these wires

	// Store the values locally in the Bob struct
	bob.e_y = e_y // for encoding (in the Encode function)
	bob.e_x = e_x // for encrypting (in the Encrypt function for OT after receiving the public keys from Alice )

	return F, d, e_xor
}

// Encoding function En uses e to map Bob's input y to a garbled input Y, by encoding each bit of y into the corresponding wire
func (bob *Bob) Encode() Y {

	// Make an array of the three bits of Bob's input y
	inputInBits := extractBits(bob.y)

	var encoded_y []string

	for i := 0; i < 3; i++ {
		// Match the input bit with the corresponding wire e_y
		encoded_y = append(encoded_y, bob.e_y[i][inputInBits[i]])
	}
	Y := Y{encoded_y, bob.e_xor}

	return Y
}

func (bob *Bob) Encrypt(elGamal *ElGamal) [][2]*Ciphertext {

	// Make 2 ciphertexts using the Public keys from Alice
	ciphertexts := make([][2]*Ciphertext, 2)

	for i := 0; i < 2; i++ {
		// Convert the string to a big.Int.
		// Assuming the string is a binary representation, you might adjust this for different representations.
		intValue1 := new(big.Int)
		intValue2 := new(big.Int)

		// Convert the strings to a big.Int
		intValue1.SetString(bob.e_x[i][0], 2) // 2 for binary
		intValue2.SetString(bob.e_x[i][1], 2) // 2 for binary

		// Encrypts both values from the wire in the same ciphertext
		ciphertexts[i][0] = elGamal.Encrypt(intValue1, bob.publicKeys[i])
		ciphertexts[i][1] = elGamal.Encrypt(intValue2, bob.publicKeys[i])
	}

	return ciphertexts

}
