package main

import (
	handin2 "Handin2"
	"fmt"
)

func main() {
	for recipient := handin2.ABplus; recipient <= handin2.Ominus; recipient++ {
		for donor := handin2.ABplus; donor <= handin2.Ominus; donor++ {

			OTTTResult := handin2.ComputeOTTTBloodTypeCompatability(recipient, donor)
			LookUPResult := handin2.LookUpBloodType(recipient, donor)

			if OTTTResult == LookUPResult {
				fmt.Println("Correct with recipient: ", handin2.GetBloodTypeName(recipient), " and donor: ", handin2.GetBloodTypeName(donor))
			} else {
				fmt.Println("Incorrect with recipient: ", handin2.GetBloodTypeName(recipient), " and donor: ", handin2.GetBloodTypeName(donor))
			}

		}
	}
}
