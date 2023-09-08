package handin2

import "math/rand"

type Dealer struct {
	r        int
	s        int
	matrix_a [8][8]bool
	matrix_b [8][8]bool
}

func (d *Dealer) Init() {
	d.r = rand.Intn(8) // generate random number between 0 and 7 for i coordinate
	d.s = rand.Intn(8) // generate random number between 0 and 7 for j coordinate

	d.matrix_b = GenerateRandomMatrix() // Generate random 8x8 matrix for Bob

	shiftMatrix := ShiftMatrix(d.r, d.s) // Generate shift matrix from r and s, and bloodtype compatibility lookup table

	d.matrix_a = XORMatrix(d.matrix_b, shiftMatrix) // Generate Alice's matrix by XOR'ing Bob's matrix with the shift matrix
}

// Getter method for matrix_a
func (d *Dealer) GetMatrixA() [8][8]bool {
	return d.matrix_a
}

// Getter method for matrix_b
func (d *Dealer) GetMatrixB() [8][8]bool {
	return d.matrix_b
}

// Getter method for r
func (d *Dealer) GetR() int {
	return d.r
}

// Getter method for s
func (d *Dealer) GetS() int {
	return d.s
}
