package handin6

import "math/big"

type Bob struct {
	y              [3]int     // Bob's input in bits (y1, y2, y3)
	encryptedYBits []*big.Int // Bob's encrypted input bits
	encryptedOnes  []*big.Int // Bob's encrypted 1's
}

// Init sets Bob's input y to be the bits y1, y2, y3
func (bob *Bob) Init(y int) {

	bob.y = ExtractBits(y)
}

// Encrypt encrypts Bob's input bits y1, y2, y3 using DHE (Stores locally for bob)
func (bob *Bob) Encrypt(DHE *DHE) {

	// Make list for Bob's encrypted input bits and encrypt the bits one at a time using DHE
	bob.encryptedYBits = make([]*big.Int, 3) // Stores locally for bob
	for i := 0; i < 3; i++ {
		bob.encryptedYBits[i] = DHE.Encrypt(bob.y[i])
	}

	// Make list of six 1's for the boolean formula
	bob.encryptedOnes = make([]*big.Int, 6) // Stores locally for bob
	for i := 0; i < 6; i++ {
		one := 1
		bob.encryptedOnes[i] = DHE.Encrypt(one)
	}

}

// Evaluate the circuit using the boolean formula:
// (1 ^ ((1 ^ x1) & y1)) & (1 ^ ((1 ^ x2) & y2)) & (1 ^ ((1 ^ x3) & y3))
func (bob *Bob) RecieveAndEvaluate(encryptedXBits []*big.Int) *big.Int {

	level3 := make([]*big.Int, 3)
	level2 := make([]*big.Int, 3)
	level1 := make([]*big.Int, 3)

	// Block 1: (1 ^ ((1 ^ x1) & y1))
	level3[0] = HE_xor(bob.encryptedOnes[0], encryptedXBits[0]) // 1 ^ x1
	level2[0] = HE_and(level3[0], bob.encryptedYBits[0])        // (1 ^ x1) & y1
	level1[0] = HE_xor(bob.encryptedOnes[1], level2[0])         // 1 ^ ((1 ^ x1) & y1)

	// Block 2: (1 ^ ((1 ^ x2) & y2))
	level3[1] = HE_xor(bob.encryptedOnes[2], encryptedXBits[1]) // 1 ^ x2
	level2[1] = HE_and(level3[1], bob.encryptedYBits[1])        // (1 ^ x2) & y2
	level1[1] = HE_xor(bob.encryptedOnes[3], level2[1])         // 1 ^ ((1 ^ x2) & y2)

	// Block 3: (1 ^ ((1 ^ x3) & y3))
	level3[2] = HE_xor(bob.encryptedOnes[4], encryptedXBits[2]) // 1 ^ x3
	level2[2] = HE_and(level3[2], bob.encryptedYBits[2])        // (1 ^ x3) & y3
	level1[2] = HE_xor(bob.encryptedOnes[5], level2[2])         // 1 ^ ((1 ^ x3) & y3)

	// AND block 1 and block 2
	levelAnd1 := HE_and(level1[0], level1[1])

	// AND block 3 and the result of AND block 1 and block 2
	levelAnd2 := HE_and(level1[2], levelAnd1)

	return levelAnd2

}
