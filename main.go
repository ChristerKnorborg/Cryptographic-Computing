package main

import handin3 "Handin3"

func main() {

	// Try all combinations of blood types and compare with the lookup table
	for recipient := handin3.Ominus; recipient <= handin3.ABplus; recipient++ {
		for donor := handin3.Ominus; donor <= handin3.ABplus; donor++ {
			aliceBloodType := recipient
			bobBloodType := donor
			BeDOZaResult := handin3.ComputeBeDOZaBloodTypeCompatability(aliceBloodType, bobBloodType)
			lookupTableResult := handin3.LookUpBloodType(aliceBloodType, bobBloodType)
			if BeDOZaResult != lookupTableResult {
				println("Error: Blood type compatibility lookup table and BeDOZa protocol does not agree for recipient: " + handin3.GetBloodTypeName(aliceBloodType) + " and donor: " + handin3.GetBloodTypeName(bobBloodType))

			}
		}
	}

}
