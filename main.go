package main

import (
	"fmt"
	h6 "handin6"
)

func main() {
	// // Try all combinations of blood types and compare with the lookup table result
	for recipient := h6.Ominus; recipient <= h6.ABplus; recipient++ {
		for donor := h6.Ominus; donor <= h6.ABplus; donor++ {
			aliceBloodType := recipient
			bobBloodType := donor
			ObliviousTransferResult := h6.HomomorphicBloodtypeEncryption(aliceBloodType, bobBloodType)
			lookupTableResult := h6.LookUpBloodType(aliceBloodType, bobBloodType)
			if ObliviousTransferResult != lookupTableResult {
				fmt.Printf("Incorrect result for recipient: %s and donor: %s\n", h6.GetBloodTypeName(recipient), h6.GetBloodTypeName(donor))
			} else {
				fmt.Printf("Correct result for recipient: %s and donor: %s\n", h6.GetBloodTypeName(recipient), h6.GetBloodTypeName(donor))
			}
		}
	}
}
