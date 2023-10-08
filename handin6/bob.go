package handin6

import "math/big"

type Bob struct {
	y              int        // Bob's input
	encryptedYBits []*big.Int // Bob's encrypted input bits
	encryptedOnes  []*big.Int // Bob's encrypted 1's
}

func (bob *Bob) Init(y int, DHE *DHE) {

	bob.y = y
	inputYInBits := ExtractBits(y)

	// Make list for Bob's encrypted input bits and encrypt the bits one at a time using DHE
	bob.encryptedYBits = make([]*big.Int, 3) // Stores locally for bob
	for i := 0; i < 3; i++ {
		bob.encryptedYBits[i] = DHE.Encrypt(inputYInBits[i])
	}

	// Make list of 6 1's for the boolean formula
	bob.encryptedOnes = make([]*big.Int, 6) // Stores locally for bob
	for i := 0; i < 6; i++ {
		one := 1
		bob.encryptedOnes[i] = DHE.Encrypt(one)

	}

}

// (1 ^ ((1 ^ x1) & y1)) & (1 ^ ((1 ^ x2) & y2)) & (1 ^ ((1 ^ x3) & y3))
func (bob *Bob) RecieveAndEvaluate(encryptedXBits []*big.Int) *big.Int {

	level3 := make([]*big.Int, 3)
	level2 := make([]*big.Int, 3)
	level1 := make([]*big.Int, 3)

	// Block 1: (1 ^ ((1 ^ x1) & y1))
	level3[0] = HE_xor(bob.encryptedOnes[0], encryptedXBits[0])
	level2[0] = HE_and(level3[0], bob.encryptedYBits[0])
	level1[0] = HE_xor(bob.encryptedOnes[1], level2[0])

	// Block 2: (1 ^ ((1 ^ x2) & y2))
	level3[1] = HE_xor(bob.encryptedOnes[2], encryptedXBits[1])
	level2[1] = HE_and(level3[1], bob.encryptedYBits[1])
	level1[1] = HE_xor(bob.encryptedOnes[3], level2[1])

	// Block 3: (1 ^ ((1 ^ x3) & y3))
	level3[2] = HE_xor(bob.encryptedOnes[4], encryptedXBits[2])
	level2[2] = HE_and(level3[2], bob.encryptedYBits[2])
	level1[2] = HE_xor(bob.encryptedOnes[5], level2[2])

	// AND block 1 and block 2
	levelAnd1 := HE_and(level1[0], level1[1])

	// AND block 3 and the result of AND block 1 and block 2
	levelAnd2 := HE_and(level1[2], levelAnd1)

	return levelAnd2

}
