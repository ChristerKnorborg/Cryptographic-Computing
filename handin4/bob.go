package handin4

import "math/big"

type Bob struct {
	y int // Bob input
}

func (bob *Bob) Init(y int) {
	// Set Bob's input as the y provided by the ObliviousTransfer function
	bob.y = y

}

func (bob *Bob) Transfer(publicKeys []*big.Int, elGamal *ElGamal) []*Ciphertext {

	// Generate m1, m2, ... , m8 from lookup table and Bob's input.
	messages := make([]*big.Int, 8)
	for i := 0; i < 8; i++ {
		messages[i] = big.NewInt(int64(Bloodtype_compatibility[i][bob.y]))
	}

	// Encrypt messages using the public keys provided by Alice
	ciphertexts := make([](*Ciphertext), 8)
	for i := 0; i < 8; i++ {
		ciphertexts[i] = elGamal.Encrypt(messages[i], publicKeys[i])
	}
	return ciphertexts
}
