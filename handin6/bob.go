package handin6

import "math/big"

type Bob struct {
	y              [3]int      // Bob's input in bits (y1, y2, y3)
	encryptedYBits [6]*big.Int // Bob's encrypted input bits
	encryptedOnes  [6]*big.Int // Bob's encrypted 1's
}

// Init sets Bob's input y to be the bits y1, y2, y3
func (bob *Bob) Init(y int) {

	bob.y = ExtractBits(y)
}

// Encrypt encrypts Bob's input bits y1, y2, y3 using DHE (Stores locally for bob)
func (bob *Bob) Encrypt(DHE *DHE) {

	// Encrypt Bob's input bits one at a time using DHE
	for i := 0; i < 3; i++ {
		bob.encryptedYBits[i] = DHE.Encrypt(bob.y[i])
	}

	// Make of six encrypted 1's for the boolean formula
	for i := 0; i < 6; i++ {
		one := 1
		bob.encryptedOnes[i] = DHE.Encrypt(one)
	}

}

// Evaluate the circuit using the boolean formula:
// (1 ^ ((1 ^ x1) & y1)) & (1 ^ ((1 ^ x2) & y2)) & (1 ^ ((1 ^ x3) & y3))
func (bob *Bob) RecieveAndEvaluate(encryptedXBits []*big.Int) *big.Int {

	depth3 := [3]*big.Int{}
	depth2 := [6]*big.Int{}

	depth3[0] = HE_xor(bob.encryptedOnes[0], encryptedXBits[0]) // 1 ^ x1
	depth3[1] = HE_xor(bob.encryptedOnes[1], encryptedXBits[1]) // 1 ^ x2
	depth3[2] = HE_xor(bob.encryptedOnes[2], encryptedXBits[2]) // 1 ^ x3

	depth2[0] = HE_and(depth3[0], bob.encryptedYBits[0]) // (1 ^ x1) & y1
	depth2[1] = HE_and(depth3[1], bob.encryptedYBits[1]) // (1 ^ x2) & y2
	depth2[2] = HE_and(depth3[2], bob.encryptedYBits[2]) // (1 ^ x3) & y3
	depth2[3] = HE_xor(depth2[0], bob.encryptedOnes[3])  // 1 ^ ((1 ^ x1) & y1)
	depth2[4] = HE_xor(depth2[1], bob.encryptedOnes[4])  // 1 ^ ((1 ^ x2) & y2)
	depth2[5] = HE_xor(depth2[2], bob.encryptedOnes[5])  // 1 ^ ((1 ^ x3) & y3)

	depth1 := HE_and(depth2[3], depth2[4]) // (1 ^ ((1 ^ x1) & y1)) & (1 ^ ((1 ^ x2) & y2))

	depth0 := HE_and(depth1, depth2[5]) // (1 ^ ((1 ^ x1) & y1)) & (1 ^ ((1 ^ x2) & y2)) & (1 ^ ((1 ^ x3) & y3))

	return depth0

}
