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

func CreateGarbledCircuit() {
	F, d, e_y, e_x, e_xor := MakeGarbledCircuit()

	e_y = Encode(e_y, bob.y)
	Y := [][][2]string{e_y, e_xor}

	return F, Y, d

}
