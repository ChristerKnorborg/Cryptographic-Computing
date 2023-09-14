package main

import handin3 "Handin3"

func main() {

	// test for all Bloodtype combinations that the XOR of shares is the same as the input bits
	for recipient := handin3.ABplus; recipient <= handin3.Ominus; recipient++ {
		for donor := handin3.ABplus; donor <= handin3.Ominus; donor++ {

			aliceBloodType := recipient
			bobBloodType := donor
			// print("Alice blood type: " + string(int(aliceBloodType)) + "\n")
			// print("Bob blood type: " + string(int(bobBloodType)) + "\n")

			alice := handin3.Alice{}
			bob := handin3.Bob{}

			// extract bits from Alice (recipient)
			x := int(aliceBloodType)
			x1 := (x >> 2) & 1 // extract 3rd rightmost bit
			x2 := (x >> 1) & 1 // extract 2nd rightmost bit
			x3 := x & 1        // extract rightmost bit

			// extract bits from Bob (donor)
			y := int(bobBloodType)
			y1 := (y >> 2) & 1 // extract 3rd rightmost bit
			y2 := (y >> 1) & 1 // extract 2nd rightmost bit
			y3 := y & 1        // extract rightmost bit

			// println(x1, x2, x3)
			// println(y1, y2, y3)

			x1_b, x2_b, x3_b := alice.TakeInput(x1, x2, x3)
			x1_a, x2_a, x3_a := alice.GetXShares()
			y1_a, y2_a, y3_a := bob.TakeInput(y1, y2, y3)
			y1_b, y2_b, y3_b := bob.GetYShares()

			// println("shares:")
			// println(x1_a, x2_a, x3_a)
			// println(x1_b, x2_b, x3_b)

			// println(y1_a, y2_a, y3_a)
			// println(y1_b, y2_b, y3_b)

			// println("XOR of shares: ")
			// println(x1_a^x1_b, x2_a^x2_b, x3_a^x3_b)
			// println(y1_a^y1_b, y2_a^y2_b, y3_a^y3_b)

			println((x1_a^x1_b) == x1, (x2_a^x2_b) == x2, (x3_a^x3_b) == x3)
			println((y1_a^y1_b) == y1, (y2_a^y2_b) == y2, (y3_a^y3_b) == y3)
		}
	}

	// Test dealer and Alice and Bob
	dealer := handin3.Dealer{}
	alice := handin3.Alice{}
	bob := handin3.Bob{}

	dealer.Init(5) // Dealer generates two shares of u, v and w for each AND gate - one for bob and one for alice
	alice.In

}
