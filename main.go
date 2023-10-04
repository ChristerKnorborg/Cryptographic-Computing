package main

import (
	"fmt"
	h5 "handin5"
)

func main() {
	// // Try all combinations of blood types and compare with the lookup table result
	for recipient := h5.Ominus; recipient <= h5.ABplus; recipient++ {
		for donor := h5.Ominus; donor <= h5.ABplus; donor++ {
			aliceBloodType := recipient
			bobBloodType := donor
			ObliviousTransferResult := h5.GarbledCircuit(aliceBloodType, bobBloodType)
			lookupTableResult := h5.LookUpBloodType(aliceBloodType, bobBloodType)
			if ObliviousTransferResult != lookupTableResult {
				fmt.Printf("Incorrect result for recipient: %s and donor: %s\n", h5.GetBloodTypeName(recipient), h5.GetBloodTypeName(donor))
			} else {
				fmt.Printf("Correct result for recipient: %s and donor: %s\n", h5.GetBloodTypeName(recipient), h5.GetBloodTypeName(donor))
			}
		}
	}
}
