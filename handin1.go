package main

import "fmt"

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

func main() {
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
}
