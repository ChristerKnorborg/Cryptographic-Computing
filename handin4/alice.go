package handin4

import (
	"fmt"
	"math/big"
)

type Alice struct {
	x  int      // Alice input
	sk *big.Int // secret key
}

func (alice *Alice) Init(x int) {
	// Set alice's input as the x provided by the ObliviousTransfer function
	alice.x = x
}

func (alice *Alice) Choose(elGamal *ElGamal) []*big.Int {

	// Generate a secret key
	alice.sk = elGamal.makeSecretKey()

	// Generate Secret Key and real public key
	publicKeys := make([]*big.Int, 8)

	for i := 0; i < 8; i++ {
		if i == alice.x {
			publicKeys[i] = elGamal.Gen(alice.sk)
		} else {
			publicKeys[i] = elGamal.OGen()
		}
	}
	return publicKeys
}

func (alice *Alice) Retrieve(ciphertexts []*Ciphertext, elGamal *ElGamal) *big.Int {
	ciphertext := ciphertexts[alice.x] // Extract the ciphertext corresponding to Alice's input x from Bob's list of ciphertexts
	c1, c2 := ciphertext.c1, ciphertext.c2

	fmt.Printf("c1: %v\n", c1)
	fmt.Printf("c2: %v\n", c2)
	// Decrypt the ciphertext
	result := elGamal.Decrypt(c1, c2, alice.sk)
	return result

}
