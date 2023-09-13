package handin3

// One-Time Truth Table (OTTT) protocol that check if recipient blood type can receive donor blood type
func ComputeOTTTBloodTypeCompatability(recipient bloodtype, donor bloodtype) bool {
	x := int(recipient) // Alice input (blood type)
	y := int(donor)     // Bob input (blood type)

	n := 5 // number of and gates in the circuit from Claudio's master solution

	d := Dealer{}
	a := Alice{}
	b := Bob{}

	d.Init(n) // Dealer generates two tuples of u, v and w for each AND gate - one for bob and one for alice

	a.Init(d.GetAliceUVW()) // Alice gets the tuples of u's, v's and w's from the dealer
	b.Init(d.GetBobUVW())   // Bob gets the tuples of u's, v's and w's from the dealer

	// Extract the bits from Alice (recipient)
	x1 := (x >> 2) & 1 // extract 3rd rightmost bit
	x2 := (x >> 1) & 1 // extract 2nd rightmost bit
	x3 := x & 1        // extract rightmost bit

	// Extract the bits from Bob (donor)
	y1 := (y >> 2) & 1 // extract 3rd rightmost bit
	y2 := (y >> 1) & 1 // extract 2nd rightmost bit
	y3 := y & 1        // extract rightmost bit

	// Layer 1: Input is negated for Alice
	z1 := x1 ^ 1
	z2 := x2 ^ 1
	z3 := x3 ^ 1

	// Layer 2: And Bob's input with Alice's negated input

	layer2()               // Bob second layer - ands the negated input with his input
	layer2(y1, y2, y3, z1) // Alice second layer - ands the negated input with his input

}
