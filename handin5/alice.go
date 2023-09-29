package handin5

import "math/big"

type Alice struct {
	x  int      // Alice input
	sk *big.Int // Secret key
}

// Set alice's input as the x provided by the GarbledCircuit function
func (alice *Alice) Init(x int) {
	alice.x = x
}

func (alice *Alice) Choose(elGamal *ElGamal) []*big.Int {

	// Generate a secret key
	alice.sk = elGamal.makeSecretKey()

	// Initialize a list of 8 public keys to be sent to Bob
	publicKeys := make([]*big.Int, 8)

	for i := 0; i < 8; i++ {
		if i == alice.x {
			publicKeys[i] = elGamal.Gen(alice.sk) // Generate the real public key corresponding to Alice's input x
		} else {
			publicKeys[i] = elGamal.OGen() // Generate 7 fake public keys using the oblivious version of Gen
		}
	}
	return publicKeys
}
