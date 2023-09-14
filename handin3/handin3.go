package handin3

func ComputeBeDOZaBloodTypeCompatability(recipient Bloodtype, donor Bloodtype) bool {
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
	   Alice also negates the output of the AND gates (XOR with constant 1).
	   Finally, they prepare the next AND between ???? */

	d_a, e_a := alice.Stage2(d1_b, d2_b, d3_b, e1_b, e2_b, e3_b)
	d_b, e_b := bob.Stage2(d1_a, d2_a, d3_a, e1_a, e2_a, e3_a)

	/* Stage 3: Alice and Bob computes the AND of z1 and z2 of the 3 outputs from the the previous layer.
	   Notice, the the final AND with z3 is the succeding stage. */

	// Alice and Bob receive masked values from their counterpart and prepare the next AND between the result of this AND and z3
	d_a2, e_a2 := alice.Stage3(d_b, e_b)
	d_b2, e_b2 := bob.Stage3(d_a, e_a)

	/* Stage 4: Alice and Bob computes the AND of last of the 3 outputs from the Stage2 and the output from the Stage3 */

	z_a := alice.Stage4(d_b2, e_b2)
	z_b := bob.Stage4(d_a2, e_a2)

	z := z_a ^ z_b // Alice and Bob XOR their shares of the output to get the final output

	return z == 1 // return true if z is 1, else return false
}

func Debug() bool {

	n := 5 // number of and gates in the circuit from Claudio's master solution

	alice := Alice{}
	bob := Bob{}

	UVWAlice := []UVW{}
	UVWBob := []UVW{}

	for i := 0; i < n; i++ {
		u, v, w := 0, 1, 1
		UVWAlice = append(UVWAlice, UVW{u, v, w})
		UVWBob = append(UVWBob, UVW{u, v, w})
	}

	alice.Init(UVWAlice) // Alice gets the shares of u's, v's and w's from the dealer
	bob.Init(UVWBob)     // Bob gets the shares of u's, v's and w's from the dealer

	// Alice and Bob generate shares of their input bits - X and Y respectively
	// x1_b, x2_b, x3_b := alice.TakeInput(x1, x2, x3)
	// y1_a, y2_a, y3_a := bob.TakeInput(y1, y2, y3)

	//x1, x2, x3 := 1, 1, 1 // Alice input (blood type)
	//y1, y2, y3 := 1, 1, 1 // Bob input (blood type)
	x1_a, x2_a, x3_a := 0, 0, 0
	x1_b, x2_b, x3_b := 1, 1, 1
	y1_a, y2_a, y3_a := 0, 0, 0
	y1_b, y2_b, y3_b := 1, 1, 1

	alice.x1, alice.x2, alice.x3 = x1_a, x2_a, x3_a
	bob.y1, bob.y2, bob.y3 = y1_b, y2_b, y3_b

	// Alice and Bob receive shares of their counterpart's input bits
	alice.ReceiveInputShares(y1_a, y2_a, y3_a)
	bob.ReceiveInputShares(x1_b, x2_b, x3_b)

	/* Stage 1: Input is negated for Alice (XOR with constant 1).
	   Alice and Bob both mask their respective input shares s.t. d = x ⊕ u, and e = y ⊕ v.
	   Since we have 3 ANDS (one for each bit) we get 3 values of each */
	d1_a, d2_a, d3_a, e1_a, e2_a, e3_a := alice.Stage1()
	d1_b, d2_b, d3_b, e1_b, e2_b, e3_b := bob.Stage1()

	// println("Stage1 Our:")
	// println("d1_a, d2_a, d3_a, e1_a, e2_a, e3_a", d1_a, d2_a, d3_a, e1_a, e2_a, e3_a)
	// println("d1_b, d2_b, d3_b, e1_b, e2_b, e3_b", d1_b, d2_b, d3_b, e1_b, e2_b, e3_b)

	/* Stage 2: Alice and Bob receive masked values from their counterpart.
	   Then, they both compute the output shares [Z] of the 3 AND gates.
	   Alice also negates the output of the AND gates (XOR with constant 1).
	   Finally, they prepare the next AND between ???? */

	d_a, e_a := alice.Stage2(d1_b, d2_b, d3_b, e1_b, e2_b, e3_b)
	d_b, e_b := bob.Stage2(d1_a, d2_a, d3_a, e1_a, e2_a, e3_a)
	println("Phase2 ours:")
	println("d_a, e_a", d_a, e_a)
	println("d_b, e_b", d_b, e_b)

	/* Stage 3: Alice and Bob computes the AND of z1 and z2 of the 3 outputs from the the previous layer.
	   Notice, the the final AND with z3 is the succeding stage. */

	// Alice and Bob receive masked values from their counterpart and prepare the next AND between the result of this AND and z3
	d_a2, e_a2 := alice.Stage3(d_b, e_b)
	d_b2, e_b2 := bob.Stage3(d_a, e_a)

	/* Stage 4: Alice and Bob computes the AND of last of the 3 outputs from the Stage2 and the output from the Stage3 */

	z_a := alice.Stage4(d_b2, e_b2)
	z_b := bob.Stage4(d_a2, e_a2)

	z := z_a ^ z_b // Alice and Bob XOR their shares of the output to get the final output

	return z == 1 // return true if z is 1, else return false
}
