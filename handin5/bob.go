package handin5

import (
	"math/big"
)

type Bob struct {
	y      int           // Bob input
	OTKeys OTPublicKeys  // Public keys from Alice
	F      []GarbledGate // Garbled circuit is a list of garbled gates
	e_x    []KeyPair     // Wires alice input bits x1, x2, x3 goes. (These are encrypt with the public keys from Alice in OT)
	e_y    []KeyPair     // Wires bob input bits y1, y2, y3 goes. (These are used in the encoding function)
}

// Set Bob's input as the y provided by the GarbledCircuit function.
// This function is simply there to simulate Bob making his own input in GarbledCircuit method.
func (bob *Bob) Init(y int) {
	bob.y = y

	// Initialize struct slices
	bob.e_x = make([]KeyPair, 0)
	bob.e_y = make([]KeyPair, 0)
	bob.F = make([]GarbledGate, 0)
	bob.OTKeys = OTPublicKeys{}
}

// Receive the public keys from Alice (used in the 1 over 2 obliviousTransfer protocol).
func (bob *Bob) ReceiveKeys(OTKeys OTPublicKeys) {
	bob.OTKeys = OTKeys
}

// Create a garbled circuit from the bloodtype compatibility formula (Claudio's Master solution from handin 1).
// For each wire i ∈ [1..T] in the circuit: choose two random strings represented by the KeyPair struct.
// For each gate g ∈ [1..G] in the circuit: construct a garbled gate using the KeyPair structs from the corresponding wires.
func (bob *Bob) MakeGarbledCircuit() ([]GarbledGate, KeyPair, []KeyPair) {

	// Initialize every wire to two random strings
	var wires [23]KeyPair
	for i := 0; i < 23; i++ {
		wires[i] = KeyPair{Random128BitString(), Random128BitString()}
	}

	// Create the circuit F from the bloodtype compatibility formula:
	// (1 ^ ((1 ^ x1) & y1)) & (1 ^ ((1 ^ x2) & y2)) & (1 ^ ((1 ^ x3) & y3))
	// For every wire, we create a garbledGate with the corresponding input and output keys.
	// For example,tThe first line "XORGate(wires[0], wires[1], wires[2])" represent
	// the left input "wire[0]", the right input "wire[1]" and the output "wire[2]".
	var F []GarbledGate

	// Block 1: (1 ^ ((1 ^ x1) & y1))
	F = append(F, XORGate(wires[0], wires[1], wires[2])) // XOR constant 1 and x1. Result is ¬x1
	F = append(F, ANDGate(wires[2], wires[3], wires[4])) // AND ¬x1 with y1. Result is z1
	F = append(F, XORGate(wires[4], wires[5], wires[6])) // XOR z1 with constant 1

	// Block 2: (1 ^ ((1 ^ x2) & y2))
	F = append(F, XORGate(wires[7], wires[8], wires[9]))    // XOR constant 1 and x2. Result is ¬x2
	F = append(F, ANDGate(wires[9], wires[10], wires[11]))  // AND ¬x2 with y2. Result is z2
	F = append(F, XORGate(wires[11], wires[12], wires[13])) // XOR z2 with constant 1

	// Block 3: (1 ^ ((1 ^ x3) & y3))
	F = append(F, XORGate(wires[14], wires[15], wires[16])) // XOR constant 1 and x3. Result is ¬x3
	F = append(F, ANDGate(wires[16], wires[17], wires[18])) // AND ¬x3 with y3. Result is z3
	F = append(F, XORGate(wires[18], wires[19], wires[20])) // XOR z3 with constant 1

	// AND block 1 and block 2
	F = append(F, ANDGate(wires[6], wires[13], wires[21])) // AND z1 and z2. Result is z4

	// AND block 3 and the result from the AND block 1 and block 2
	F = append(F, ANDGate(wires[20], wires[21], wires[22])) // AND z3 and z4. Final

	// Define d = (Z_0, Z_1) = (K^T_0 , K^T_1)
	d := wires[22]

	e_x := []KeyPair{wires[1], wires[8], wires[15]}                                   // Alice input bits x1, x2, x3 goes in these wires respectively
	e_y := []KeyPair{wires[3], wires[10], wires[17]}                                  // Bob input bits y1, y2, y3 goes in these wires respectively
	e_xor := []KeyPair{wires[0], wires[5], wires[7], wires[12], wires[14], wires[19]} // constants from XOR gates goes in these wires

	// Store the values locally in the Bob struct (Since he need them later)
	bob.e_y = e_y // for encoding (in the Encode function)
	bob.e_x = e_x // for encrypting (in the Encrypt function for OT after receiving the public keys from Alice)

	return F, d, e_xor
}

// Encoding function En uses e to map Bob's input y to a garbled input Y, by encoding each bit of y into the corresponding wire
func (bob *Bob) Encode() []string {

	// Make an slice of the three bits of Bob's input y
	inputInBits := ExtractBits(bob.y) // [y1, y2, y3]

	var Y []string // Single key for each wire

	for i := 0; i < 3; i++ {
		// Match the input bit with the corresponding wire e_y
		if inputInBits[i] == 0 {
			Y = append(Y, bob.e_y[i].K_0)
		} else if inputInBits[i] == 1 {
			Y = append(Y, bob.e_y[i].K_1)
		} else {
			panic("Input bit is not 0 or 1")
		}
	}
	return Y
}

func (bob *Bob) Encrypt(elGamal *ElGamal) [3][2]*Ciphertext {

	// Make 2 ciphertexts (of c1, c2) for each of Alice three input wires.
	// Two ciphertexts is due to encryption of both bits for each wire where one of the bits use is a real key and the other use a fake
	encrypted_x := [3][2]*Ciphertext{}

	for i := 0; i < 3; i++ {
		// Convert the strings to a big.Int to be used in the ElGamal encryption from last week (Base 16 is for hexadecimal string)
		wireK_0, err0 := new(big.Int).SetString(bob.e_x[i].K_0, 16)
		wireK_1, err1 := new(big.Int).SetString(bob.e_x[i].K_1, 16)

		if !err0 || !err1 {
			panic("Could not convert string to big.Int")
		}
		encrypted_x[i][0] = elGamal.Encrypt(wireK_0, bob.OTKeys.Keys[i][0])
		encrypted_x[i][1] = elGamal.Encrypt(wireK_1, bob.OTKeys.Keys[i][1])
	}
	return encrypted_x
}
