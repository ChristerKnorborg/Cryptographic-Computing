package handin6

import "math/big"

type Alice struct {
	x int // Alice's input
}

func (alice *Alice) Init(x int, DHE *DHE) []*big.Int {

	alice.x = x
	inputXInBits := ExtractBits(x)

	// Make list for Alice's encrypted bits and encrypt the bits one at a time using DHE
	encryptedXBits := make([]*big.Int, 3)
	for i := 0; i < 3; i++ {
		encryptedXBits[i] = DHE.Encrypt(inputXInBits[i])
	}

	return encryptedXBits
}
