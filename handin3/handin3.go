package handin3

// One-Time Truth Table (OTTT) protocol that check if recipient blood type can receive donor blood type
func ComputeBeDOZaBloodTypeCompatability(recipient bloodtype, donor bloodtype) bool {
	x := int(recipient) // Alice input (blood type)
	y := int(donor)     // Bob input (blood type)

	// Extract the bits from Alice (recipient)
	x1 := (x >> 2) & 1 // extract 3rd rightmost bit
	x2 := (x >> 1) & 1 // extract 2nd rightmost bit
	x3 := x & 1        // extract rightmost bit

	// Extract the bits from Bob (donor)
	y1 := (y >> 2) & 1 // extract 3rd rightmost bit
	y2 := (y >> 1) & 1 // extract 2nd rightmost bit
	y3 := y & 1        // extract rightmost bit

	n := 5 // number of and gates in the circuit from Claudio's master solution

	dealer := Dealer{}
	alice := Alice{}
	bob := Bob{}

	dealer.Init(n) // Dealer generates two shares of u, v and w for each AND gate - one for bob and one for alice

	alice.Init(dealer.GetAliceUVW()) // Alice gets the shares of u's, v's and w's from the dealer
	bob.Init(dealer.GetBobUVW())     // Bob gets the shares of u's, v's and w's from the dealer

	// Alice and Bob generate shares of their input bits - X and Y respectively
	x1_b, x2_b, x3_b := alice.TakeInput(x1, x2, x3)
	y1_a, y2_a, y3_a := bob.TakeInput(y1, y2, y3)

	// Alice and Bob receive shares of their counterpart's input bits
	alice.ReceiveInputShares(y1_a, y2_a, y3_a)
	bob.ReceiveInputShares(x1_b, x2_b, x3_b)

	/* Stage 1: Input is negated for Alice (XOR with constant 1).
	   Alice and Bob both mask their respective input shares s.t. d = x ⊕ u, and e = y ⊕ v.
	   Since we have 3 ANDS (one for each bit) we get 3 values of each */
	d1_a, d2_a, d3_a, e1_a, e2_a, e3_a := alice.Stage1()
	d1_b, d2_b, d3_b, e1_b, e2_b, e3_b := bob.Stage1()

	/* Stage 2: Alice and Bob receive masked values from their counterpart.
	   Then, they both compute the output shares [Z] of the 3 AND gates.
	   Finally, Alice also negates the output of the AND gates (XOR with constant 1) */

	alice.Stage2(d1_b, d2_b, d3_b, e1_b, e2_b, e3_b)
	bob.Stage2(d1_a, d2_a, d3_a, e1_a, e2_a, e3_a)

	// Stage 3: Alice and Bob computes the AND of 2 of the 3 outputs from the the previous layer (the final AND is in layer4)
	// Currently, Alice and Bob both have a share of [z], from the previous layer.

	// Alice and Bob mask their respective shares of [z] s.t. d = z ⊕ u, and e = z ⊕ v.
	e_a, d_a := alice.MaskZ1AndZ2() // Alice masks the output of the AND gates
	e_b, d_b := bob.MaskZ1AndZ2()   // Bob masks the output of the AND gates

	// Alice and Bob receive masked values from their counterpart
	alice.ReceiveMasked(e_b, d_b) // Alice receives masked values from Bob and unmasks them using her own shares of d and e
	bob.ReceiveMasked(e_a, d_a)   // Bob receives masked values from Aliceand unmasks them using his own shares of d and e

	// Alice and Bob compute the output shares.
	alice.ComputeZinAND() //MISSING
	bob.ComputeZinAND()   //MISSING

	// Layer 4: Alice and Bob computes the AND of last of the 3 outputs from the layer2 and the output from the layer3
	// Currently, Alice and Bob both have a share of [z], from the previous layer.
	e_a, d_a = alice.MaskZ1AndZ2() // Alice masks the output of the AND gates
	e_b, d_b = bob.MaskZ1AndZ2()   // Bob masks the output of the AND gates

	// Alice and Bob receive masked values from their counterpart
	alice.ReceiveMasked(e_b, d_b) // Alice receives masked values from Bob and unmasks them using her own shares of d and e

}
