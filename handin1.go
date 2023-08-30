package main

import "fmt"

type bloodtype int32

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

var bloodtype_compatibility [8][8]bool = [8][8]bool{
	{true, true, true, true, true, true, true, true},        // AB+
	{false, true, false, true, false, true, false, false},   // AB-
	{false, false, true, true, false, false, false, false},  // B+
	{false, false, false, true, false, false, false, false}, // B-
	{false, false, false, false, true, true, false, false},  // A+
	{false, false, false, false, false, true, false, false}, // A-
	{false, false, false, false, false, false, true, true},  // O+
	{false, false, false, false, false, false, false, true}, // O-
}

func LookUpBloodType(recipient bloodtype, donor bloodtype) bool {
	return bloodtype_compatibility[recipient][donor]
}

// BooleanFormula checks if blood type x can receive from y using Boolean operations.
// x and y should be 3-bit encoded blood types.
// func BooleanFormula(x uint8, y uint8) bool {
// 	// Extract individual bits from x and y
// 	x1 := (x >> 2) & 1
// 	x2 := (x >> 1) & 1
// 	x3 := x & 1

// 	y1 := (y >> 2) & 1
// 	y2 := (y >> 1) & 1
// 	y3 := y & 1

// 	// Use Boolean formula to determine compatibility
// 	condition1 := (x1) || y1
// 	condition2 := (!x2) || y2
// 	condition3 := (!x3) || y3

// 	return condition1 && condition2 && condition3
// }

func main() {
	fmt.Printf("%t\n", LookUpBloodType(ABplus, ABplus))

	fmt.Printf("%t\n", LookUpBloodType(Ominus, Ominus))
	fmt.Printf("%t\n", LookUpBloodType(Ominus, ABplus))
}
