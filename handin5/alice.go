package handin5

import "math/big"

type Alice struct {
	x   int         // Alice input
	sk  *big.Int    // El Gamal secret key
	e_x [][2]string // Encode values for Alice
	d   [][2]string // The Z values from the ouput of the garbled circuit
}

// Set alice's input as the x provided by the GarbledCircuit function
func (alice *Alice) Init(x int) {
	alice.x = x
}

func (alice *Alice) MakeAndTransferKeys(elGamal *ElGamal) []*big.Int {

	// Generate a secret key
	alice.sk = elGamal.makeSecretKey()

	// Make two public keys for the ObliviousTransfer
	publicKeys := make([]*big.Int, 2)
	publicKeys[0] = elGamal.Gen(alice.sk)
	publicKeys[1] = elGamal.OGen()

	return publicKeys
}

func (alice *Alice) EvaluateGarbledCircuit(F, d, Y) int {

	// Alice runs Ev to evaluate the garbled circuit [F] on the garbled input [X] and produces a garbled output [Z′]

	if Z == alice.d[0] {
		return 0
	} else if Z == alice.d[1] {
		return 1
	} else {
		panic("Decoding failed. Result is neither Z_0 or Z_1 from d = (Z_0, Z_1)")
	}
}
