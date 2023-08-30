package main

import "fmt"

type bloodtype int32

const (
	Ominus  bloodtype = 0
	Oplus             = 1
	Aminus            = 2
	Aplus             = 3
	Bminus            = 4
	Bplus             = 5
	ABminus           = 6
	ABplus            = 7
)

var bloodtypes [8][8]bool = [8][8]bool{
	{true, true, true, true, true, true, true, true},        // AB+
	{false, true, false, true, false, true, false, false},   // AB-
	{false, false, true, true, false, false, false, false},  // B+
	{false, false, false, true, false, false, false, false}, // B-
	{false, false, false, false, true, true, false, false},  // A+
	{false, false, false, false, false, true, false, false}, // A-
	{false, false, false, false, false, false, true, true},  // O+
	{false, false, false, false, false, false, false, true}, // O-
}

func tester(x bloodtype, y bloodtype) bool {
	return bloodtypes[x][y]
}

func main() {
	fmt.Printf("%t\n", tester(ABplus, ABplus))
	fmt.Printf("%t\n", tester(ABplus, Ominus))
}
