package handin6

import "math/big"

type Alice struct {
	x  [3]int   // Alice's input in bits (x1, x2, x3)
	sk *big.Int // Alice's secret key of p
}

// Init sets Alice's input x to be the bits x1, x2, x3
func (alice *Alice) Init(x int) {

	alice.x = ExtractBits(x)
}

// Encrypt encrypts Alice's input bits x1, x2, x3 using DHE and sends them to Bob
func (alice *Alice) Encrypt(DHE *DHE) []*big.Int {

	// Make list for Alice's encrypted input bits and encrypt the bits one at a time using DHE
	encryptedXBits := make([]*big.Int, 3)
	for i := 0; i < 3; i++ {
		encryptedXBits[i] = DHE.Encrypt(alice.x[i])
	}

	return encryptedXBits
}

// RecieveSecretKey recieves the secret key p from DHE
func (alice *Alice) RecieveSecretKey(p *big.Int) {
	alice.sk = p
}

// Decrypt decrypts the result of the circuit evaluation using DHE with Alice's secret key p
func (alice *Alice) Decrypt(evaluatedOutput *big.Int, DHE *DHE) int {

	// Decrypt the result
	decryptedOutput := DHE.Decrypt(evaluatedOutput, alice.sk)

	return decryptedOutput // 1 = true, 0 = false

}
