package handin5

import "math/big"

type Bob struct {
	y   int        // Bob input
	d   int        // decode value bob
	F   []*big.Int // Values from each gate of the circuit
	e_y *big.Int   // Input encoding Bob
}

// Set Bob's input as the y provided by the GarbledCircuit function
func (bob *Bob) Init(y int) {
	bob.y = y
}
