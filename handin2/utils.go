package handin2

import (
	"math/rand"
)

// XOR function that returns true if x and y are different, and false if they are the same
func XOR(x bool, y bool) bool {
	return (x || y) && !(x && y)
}

// XORMatrix function that returns entry-wise XOR of two 8x8 matrices
func XORMatrix(matrixX [8][8]bool, matrixY [8][8]bool) [8][8]bool {
	resultMatrix := [8][8]bool{}
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			resultMatrix[i][j] = XOR(matrixX[i][j], matrixY[i][j]) // XOR entry-wise
		}
	}
	return resultMatrix
}

func ShiftMatrix(r int, s int) [8][8]bool {
	shiftMatrix := [8][8]bool{}
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			shiftMatrix[i][j] = bloodtype_compatibility[(i-r+8)%8][(j-s+8)%8] // Go does not support negative modulo, so we add 8 to the result and take modulo 8 again
		}
	}
	return shiftMatrix
}

func GenerateRandomMatrix() [8][8]bool {
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

func GetBloodTypeName(bType bloodtype) string {
	switch bType {
	case ABplus:
		return "AB+"
	case ABminus:
		return "AB-"
	case Bplus:
		return "B+"
	case Bminus:
		return "B-"
	case Aplus:
		return "A+"
	case Aminus:
		return "A-"
	case Oplus:
		return "O+"
	case Ominus:
		return "O-"
	default:
		return "Unknown"
	}
}
