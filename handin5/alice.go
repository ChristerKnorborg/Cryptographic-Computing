package handin5

import "math/big"

type Alice struct {
	x   int          // Alice input
	sk  *big.Int     // El Gamal secret key
	e_x []GarbleKeys // Encode values for Alice
	d   GarbleKeys   // The Z values from the ouput of the garbled circuit
}

// Set alice's input as the x provided by the GarbledCircuit function
func (alice *Alice) Init(x int) {
	alice.x = x
}

func (alice *Alice) Choose(elGamal *ElGamal) []*big.Int {

	// Generate a secret key
	alice.sk = elGamal.makeSecretKey()

	// Make two public keys for the ObliviousTransfer
	publicKeys := make([]*big.Int, 2)
	publicKeys[0] = elGamal.Gen(alice.sk)
	publicKeys[1] = elGamal.OGen()

	return publicKeys
}
