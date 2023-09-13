package handin3

// One-Time Truth Table (OTTT) protocol that check if recipient blood type can receive donor blood type
func ComputeOTTTBloodTypeCompatability(recipient bloodtype, donor bloodtype) bool {
	x := int(recipient) // Alice input (blood type)
	y := int(donor)     // Bob input (blood type)

	n := 5 // number of and gates in the circuit from Claudio's master solution

	dealer := Dealer{}
	alice := Alice{}
	bob := Bob{}

	dealer.Init(n) // Dealer generates two tuples of u, v and w for each AND gate - one for bob and one for alice

	alice.Init(dealer.GetAliceUVW()) // Alice gets the tuples of u's, v's and w's from the dealer
	bob.Init(dealer.GetBobUVW())     // Bob gets the tuples of u's, v's and w's from the dealer

	// Extract the bits from Alice (recipient)
	x1 := (x >> 2) & 1 // extract 3rd rightmost bit
	x2 := (x >> 1) & 1 // extract 2nd rightmost bit
	x3 := x & 1        // extract rightmost bit

	// Extract the bits from Bob (donor)
	y1 := (y >> 2) & 1 // extract 3rd rightmost bit
	y2 := (y >> 1) & 1 // extract 2nd rightmost bit
	y3 := y & 1        // extract rightmost bit

	// Alice and Bob generate shares of their input bits
	x1_b, x2_b, x3_b := alice.TakeInput(x1, x2, x3)
	y1_a, y2_a, y3_a := bob.TakeInput(y1, y2, y3)

	// Alice and Bob receive shares of their counterpart's input bits
	alice.ReceiveInput(y1_a, y2_a, y3_a)
	bob.ReceiveInput(x1_b, x2_b, x3_b)

	// Layer 1: Input is negated for Alice (XOR with constant 1). Bob does nothing.
	alice.x1 = alice.x1 ^ 1
	alice.x2 = alice.x2 ^ 1
	alice.x3 = alice.x3 ^ 1

	// Layer 2: And Bob's input with Alice's negated input

	// Alice and Bob both mask their respective input shares s.t. d = x ⊕ u, and e = y ⊕ v.
	d1_a, d2_a, d3_a, e1_a, e2_a, e3_a := alice.MaskXandY()
	d1_b, d2_b, d3_b, e1_b, e2_b, e3_b := bob.MaskXandY()

	// Alice and Bob receive masked values from their counterpart
	alice.ReceiveValues(d1_b, d2_b, d3_b, e1_b, e2_b, e3_b)
	bob.ReceiveValues(d1_a, d2_a, d3_a, e1_a, e2_a, e3_a)

	// Alice and Bob compute the output of the 3 AND gates
	z1_a, z2_a, z3_a := alice.ComputeZinAND()
	z1_b, z2_b, z3_b := bob.ComputeZinAND()

	// Layer 3: Alice negate the output of the AND gates (XOR with constant 1). Bob does nothing.
	z1_a = z1_a ^ 1
	z2_a = z2_a ^ 1
	z3_a = z3_a ^ 1

	// Layer 4: Alice and Bob computes the AND of the previous layer's output
	alice.Mask(z1_a, z2_a, z3_a) // Alice masks the output of the AND gates
	bob.Mask(z1_b, z2_b, z3_b)   // Bob masks the output of the AND gates

	// Alice and Bob receive masked values from their counterpart
	alice.ReceiveMaskedVal2() // Alice receives masked values from Bob
	bob.ReceiveMaskedVal2()   // Bob receives masked values from Alice

}
