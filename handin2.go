package main

import (
	"fmt"
	"math/rand"
)

type Dealer struct {
	r        int
	s        int
	matrix_a [8][8]bool
	matrix_b [8][8]bool
}

type Alice struct {
	x        int
	matrix_a [8][8]bool
	u        int
	r        int
	v        int
	z_B      bool
}

type Bob struct {
	y        int
	matrix_b [8][8]bool
	s        int
	u        int
}

func (d *Dealer) init() {
	d.r = rand.Intn(8) // for i coordinate
	d.s = rand.Intn(8) // for j coordinate

	d.matrix_b = generateRandomMatrix() // Generate random 8x8 matrix for Bob

	shiftMatrix := shiftMatrix(d.r, d.s) // Generate shift matrix

	d.matrix_a = XORMatrix(d.matrix_b, shiftMatrix) // Generate Alice's matrix by XOR'ing Bob's matrix with the shift matrix
}

func (d *Dealer) getMatrixA() [8][8]bool {
	return d.matrix_a
}

func (d *Dealer) getMatrixB() [8][8]bool {
	return d.matrix_b
}

func (d *Dealer) getR() int {
	return d.r
}

func (d *Dealer) getS() int {
	return d.s
}

func (a *Alice) init(x int, matrix [8][8]bool, r int) {
	a.x = x
	a.matrix_a = matrix
	a.r = r

}

func (b *Bob) init(y int, matrix [8][8]bool, s int) {
	b.y = y
	b.matrix_b = matrix
	b.s = s
}

func (a *Alice) send() int {
	u := (int(a.x) + a.r) % 8 // Alice computes u = x + r mod 2n and sends u to Bob
	a.u = u                   // Alice stores u
	return u
}

func (a *Alice) receive(v int, z_B bool) {
	a.v = v
	a.z_B = z_B
}

func (a *Alice) computeOutput() bool {
	z := XOR(a.matrix_a[a.u][a.v], a.z_B) // Alice outputs z = M_A[u, v] ⊕ z_B u
	return z
}

func (b *Bob) send() (int, bool) {
	v := (int(b.y) + b.s) % 8 // Bob computes v = y + s mod 2n,
	z_B := b.matrix_b[v][b.s] // and z_B = M_B[u, v] and sends (v, z_B) to Alice
	return v, z_B
}

func (b *Bob) receive(u int) {
	b.u = u
}

func XOR(x bool, y bool) bool {
	return (x || y) && !(x && y)
}

func XORMatrix(matrixX [8][8]bool, matrixY [8][8]bool) [8][8]bool {
	resultMatrix := [8][8]bool{}
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			resultMatrix[i][j] = XOR(matrixX[i][j], matrixY[i][j])
		}
	}
	return resultMatrix
}

func shiftMatrix(r int, s int) [8][8]bool {
	shiftMatrix := [8][8]bool{}
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			shiftMatrix[i][j] = bloodtype_compatibility[(i-r+8)%8][(j-s+8)%8] // need to add 8 to make sure we get positive numbers. Go does not support negative modulo
		}
	}
	return shiftMatrix
}

func generateRandomMatrix() [8][8]bool {
	// rand.Seed(time.Now().UnixNano())

	matrix := [8][8]bool{}
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			if rand.Intn(2) == 1 { // generate random 0 or 1
				matrix[i][j] = true // true = 1
			} else {
				matrix[i][j] = false // false = 0
			}
		}
	}

	return matrix
}

type bloodtype uint8

const (
	ABplus  bloodtype = 0
	ABminus bloodtype = 1
	Bplus   bloodtype = 2
	Bminus  bloodtype = 3
	Aplus   bloodtype = 4
	Aminus  bloodtype = 5
	Oplus   bloodtype = 6
	Ominus  bloodtype = 7
)

// bloodtype compatibility lookup table
var bloodtype_compatibility [8][8]bool = [8][8]bool{
	{true, true, true, true, true, true, true, true},        // AB+
	{false, true, false, true, false, true, false, true},    // AB-
	{false, false, true, true, false, false, true, true},    // B+
	{false, false, false, true, false, false, false, true},  // B-
	{false, false, false, false, true, true, true, true},    // A+
	{false, false, false, false, false, true, false, true},  // A-
	{false, false, false, false, false, false, true, true},  // O+
	{false, false, false, false, false, false, false, true}, // O-
}

// LookUpBloodType checks if recipient blood type can receive donor blood type using lookup table
func LookUpBloodType(recipient bloodtype, donor bloodtype) bool {
	return bloodtype_compatibility[recipient][donor]
}

// main function for testing that the protocol works. The dealer gives Alice and Bob their random matrices and coordinates.
// Alice and Bob both initialize their matrices and coordinates. Alice sends u to Bob, Bob sends v and z_B to Alice.
func main() {
	AliceBloodType := ABplus
	BobBloodType := Ominus
	fmt.Println("Alice's blood type is", int(AliceBloodType), "and Bob's blood type is", int(BobBloodType))

	d := Dealer{}
	d.init() // Dealer initializes the matrices and coordinates

	a := Alice{}
	a.init(int(AliceBloodType), d.getMatrixA(), d.getR()) // Alice initializes her matrix and coordinates

	b := Bob{}
	b.init(int(BobBloodType), d.getMatrixB(), d.getS()) // Bob initializes his matrix and coordinates

	b.receive(a.send()) // Alice sends u to Bob

	v, z_B := b.send()
	a.receive(v, z_B) // Bob sends v and z_B to Alice

	fmt.Println(a.computeOutput())
}
