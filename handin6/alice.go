package handin6

import "math/big"

type Alice struct {
	x  [3]int   // Alice's input in bits (x1, x2, x3)
	sk *big.Int // Alice's secret key of p
}

func (alice *Alice) Init(x int) {

	alice.x = ExtractBits(x)
}

func (alice *Alice) Encrypt(DHE *DHE) []*big.Int {

	// Make list for Alice's encrypted input bits and encrypt the bits one at a time using DHE
	encryptedXBits := make([]*big.Int, 3) // Stores locally for alice
	for i := 0; i < 3; i++ {
		encryptedXBits[i] = DHE.Encrypt(alice.x[i])
	}

	return encryptedXBits
}

func (alice *Alice) RecieveSecretKey(p *big.Int) {
	alice.sk = p
}

func (alice *Alice) Decrypt(evaluatedOutput *big.Int, DHE *DHE) int {

	// Decrypt the result
	decryptedOutput := DHE.Decrypt(evaluatedOutput, alice.sk)

	return decryptedOutput

}
