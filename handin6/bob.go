package handin6

import "math/big"

type Bob struct {
	y int // Bob's input
}

func (bob *Bob) Init(y int, DHE *DHE) ([]*big.Int, []*big.Int) {

	bob.y = y
	inputYInBits := ExtractBits(y)

	// Make list for Bob's encrypted input bits and encrypt the bits one at a time using DHE
	encryptedYBits := make([]*big.Int, 3)
	for i := 0; i < 3; i++ {
		encryptedYBits[i] = DHE.Encrypt(inputYInBits[i])
	}

	// Make list of 6 1's for the boolean formula
	encryptedOnes := make([]*big.Int, 6)
	for i := 0; i < 6; i++ {
		one := 1
		encryptedOnes[i] = DHE.Encrypt(one)

	}

	return encryptedYBits, encryptedOnes

}
