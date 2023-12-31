package handin4

import (
	"math/big"
)

type Alice struct {
	x  int      // Alice input
	sk *big.Int // Secret key
}

func (alice *Alice) Init(x int) {
	// Set alice's input as the x provided by the ObliviousTransfer function
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

func (alice *Alice) Retrieve(ciphertexts []*Ciphertext, elGamal *ElGamal) *big.Int {
	ciphertext := ciphertexts[alice.x] // Extract the ciphertext corresponding to Alice's input x from Bob's list of ciphertexts
	c1, c2 := ciphertext.c1, ciphertext.c2

	result := elGamal.Decrypt(c1, c2, alice.sk) // Decrypt the ciphertext using Alice's secret key
	return result

}
