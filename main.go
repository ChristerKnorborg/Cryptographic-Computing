package main

import (
	h4 "Handin4"
	"fmt"
)

func main() {

	// // Try all combinations of blood types and compare with the lookup table
	for recipient := h4.Ominus; recipient <= h4.ABplus; recipient++ {
		for donor := h4.Ominus; donor <= h4.ABplus; donor++ {
			aliceBloodType := recipient
			bobBloodType := donor
			ObliviousTransferResult := h4.ObliviousTransfer(aliceBloodType, bobBloodType)
			lookupTableResult := h4.LookUpBloodType(aliceBloodType, bobBloodType)
			if ObliviousTransferResult != lookupTableResult {
				fmt.Printf("Incorrect result for recipient: %s and donor: %s\n", h4.GetBloodTypeName(recipient), h4.GetBloodTypeName(donor))
			} else {
				fmt.Printf("Correct result for recipient: %s and donor: %s\n", h4.GetBloodTypeName(recipient), h4.GetBloodTypeName(donor))
			}

		}
	}

	//Try a single combination of blood types and compare with the lookup table
	// aliceBloodType := h4.Ominus
	// bobBloodType := h4.Ominus
	// ObliviousTransferResult := h4.ObliviousTransfer(aliceBloodType, bobBloodType)
	// lookupTableResult := h4.LookUpBloodType(aliceBloodType, bobBloodType)
	// print("ObliviousTransferResult: ")
	// println(ObliviousTransferResult)
	// print("lookupTableResult: ")
	// println(lookupTableResult)

}
