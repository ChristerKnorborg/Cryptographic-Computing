package main

import (
	"math/rand"
	"time"
)

type Dealer struct {
	r        int
	s        int
	matrix_a [8][8]bool
	matrix_b [8][8]bool
}

type Alice struct {
	x        bloodtype
	matrix_a [8][8]bool
	r        int
}

type Bob struct {
	y        bloodtype
	matrix_b [8][8]bool
	s        int
	u        int
}

func (d *Dealer) init() {
	d.r = rand.Intn(7) // for i coordinate
	d.s = rand.Intn(7) // for j coordinate

	d.matrix_b = generateRandomMatrix()

	shiftMatrix := shiftMatrix(d.r, d.s)

	d.matrix_a = XORMatrix(d.matrix_b, shiftMatrix)
}

func (d *Dealer) getRandA() ([8][8]bool, int) {
	return d.matrix_a, d.r
}

func (d *Dealer) getRandB() ([8][8]bool, int) {
	return d.matrix_b, d.s
}

func (a *Alice) init(x bloodtype, d Dealer) {
	matrix, r := d.getRandA()
	a.matrix_a = matrix
	a.r = r

}

func (b *Bob) init(y bloodtype, d Dealer) {
	matrix, s := d.getRandB()
	b.matrix_b = matrix
	b.s = s
}

func (a *Alice) send() int {
	u := (a.x + a.r) % 8
	return u
}

func (a *Alice) receive() {

}

func (b *Bob) send() int {
	v := (b.y + b.s) % 8
	z_b := b.matrix_b[v][b.s]
	return v
}

func (b *Bob) receive(Alice Alice) {
	b.u = Alice.send()
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
			shiftMatrix[i][j] = bloodtype_compatibility[i-r%8][j-s%8]
		}
	}

	return shiftMatrix

}

func generateRandomMatrix() [8][8]bool {
	rand.Seed(time.Now().UnixNano())

	matrix := [8][8]bool{}
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			matrix[i][j] = rand.Intn(2) == 1
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

// BooleanFormula checks if recipient blood type can receive donor blood type using Boolean formulation
func BooleanFormula(recipient bloodtype, donor bloodtype) bool {

	x1 := (recipient >> 2) & 1 // extract 3rd rightmost bit
	x2 := (recipient >> 1) & 1 // extract 2nd rightmost bit
	x3 := recipient & 1        // extract rightmost bit

	y1 := (donor >> 2) & 1 // extract 3rd rightmost bit
	y2 := (donor >> 1) & 1 // extract 2nd rightmost bit
	y3 := donor & 1        // extract rightmost bit

	condition1 := (x1 == 0) || (y1 == 1)
	condition2 := (x2 == 0) || (y2 == 1)
	condition3 := (x3 == 0) || (y3 == 1)

	return condition1 && condition2 && condition3

}

/* func main() {
	for recipient := ABplus; recipient <= Ominus; recipient++ {
		for donor := ABplus; donor <= Ominus; donor++ {
			lookupResult := LookUpBloodType(recipient, donor)
			formulaResult := BooleanFormula(recipient, donor)
			if lookupResult != formulaResult {
				fmt.Printf("Mismatch! Recipient: %d, Donor: %d, Lookup: %t, Formula: %t\n", recipient, donor, lookupResult, formulaResult)
			} else {
				fmt.Printf("Match! Recipient: %d, Donor: %d, Result: %t\n", recipient, donor, lookupResult)
			}
		}
	}
} */
