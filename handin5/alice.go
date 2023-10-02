package handin5

import "math/big"

type Alice struct {
	x   int        // Alice input
	sk  []*big.Int // El Gamal secret keys - one to encrypt each of Alice's three input bits
	e_x []KeyPair  // Encode values for Alice
	d   []KeyPair  // The Z values from the ouput of the garbled circuit
}

// Set alice's input as the x provided by the GarbledCircuit function
func (alice *Alice) Init(x int) {
	alice.x = x
}

type OTPublicKeys struct {
	// Public keys from Alice used in the 1 over 2 obliviousTransfer protocol.
	// One key pair for each of Alice's three input bits. Each pair contains a real and a fake key,
	// where the real key corresponds to the chosen input bit of Alice.
	keys [3][2]*big.Int
}

func (alice *Alice) MakeAndTransferKeys(elGamal *ElGamal) OTPublicKeys {

	// Generate three secret keys for Alice - one for each of her three input bits
	for i := 0; i < 3; i++ {
		alice.sk[i] = elGamal.makeSecretKey()
	}

	// Extract the bits from Alice's input
	inputInBits := ExtractBits(alice.x) // [x1, x2, x3]

	OTKeys := OTPublicKeys{}

	// Encode Alice's input bits. If input bit is 0, then the real key is the first key in the pair. Otherwise, opposite.
	for i := 0; i < 3; i++ {
		if inputInBits[i] == 0 {
			OTKeys.keys[i][0] = elGamal.Gen(alice.sk[i]) // Real key
			OTKeys.keys[i][1] = elGamal.OGen()           // Fake key
		} else {
			OTKeys.keys[i][0] = elGamal.OGen()           // Fake key
			OTKeys.keys[i][1] = elGamal.Gen(alice.sk[i]) // Real key
		}
	}
	return OTKeys
}

func (alice *Alice) EvaluateGarbledCircuit(F [][]string, e_x [][2]string, d [2]string, Y Y, e_xor [][2]string) int {

	// Block 1: x1 and y1
	notX1 := EvaluateGarbledGate(F[0], alice.e_x[0], e_xor[0][1]) // XOR constant 1 and x1. Result is ¬x1
	z1 := EvaluateGarbledGate(F[1], notX1, Y.encoded_y[0])        // AND ¬x1 with y1. Result is z1
	notZ1 := EvaluateGarbledGate(F[2], z1, e_xor[1][1])           // XOR z1 with constant 1

	// Block 2: x2 and y2
	notX2 := EvaluateGarbledGate(F[3], alice.e_x[1], e_xor[2][1]) // XOR constant 1 and x2. Result is ¬x2
	z2 := EvaluateGarbledGate(F[4], notX2, Y.encoded_y[1])        // AND ¬x2 with y2. Result is z2
	notZ2 := EvaluateGarbledGate(F[5], z2, e_xor[3][1])           // XOR z2 with constant 1

	// Block 3: x3 and y3
	notX3 := EvaluateGarbledGate(F[6], alice.e_x[2], e_xor[4][1]) // XOR constant 1 and x3. Result is ¬x3
	z3 := EvaluateGarbledGate(F[7], notX3, Y.encoded_y[2])        // AND ¬x3 with y3. Result is z3
	notZ3 := EvaluateGarbledGate(F[8], z3, e_xor[5][1])           // XOR z3 with constant 1

	z4 := EvaluateGarbledGate(F[9], notZ1, notZ2) // AND ¬z1 and ¬z2. Result is z4

	Z := EvaluateGarbledGate(F[10], notZ3, z4) // AND ¬z3 and z4. Final

	// Alice runs Ev to evaluate the garbled circuit [F] on the garbled input [X] and produces a garbled output [Z′]
	if Z == d[0] {
		return 0
	} else if Z == d[1] {
		return 1
	} else {
		panic("Decoding failed. Result is neither Z_0 or Z_1 from d = (Z_0, Z_1)")
	}
}

func (alice *Alice) Decrypt(ciphertexts [][2]*Ciphertext, elGamal *ElGamal) string {
	var result string

	// Loop through the ciphertexts
	for i, _ := range ciphertexts {
		// Decrypt each value in the pair
		decryptedValue1 := elGamal.Decrypt(ciphertexts[i], alice.sk)
		decryptedValue2 := elGamal.Decrypt(pair[1], alice.sk)

		// Convert the decrypted big.Int values to binary strings
		binaryStr1 := decryptedValue1.Text(2) // 2 for binary
		binaryStr2 := decryptedValue2.Text(2) // 2 for binary

		// Append or process the binary strings as required (e.g., concatenate)
		result += binaryStr1 + binaryStr2
	}

	return result
}