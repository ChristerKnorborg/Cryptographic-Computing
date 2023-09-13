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

	// Layer 1: Input is negated for Alice
	alice.x1 = alice.x1 ^ 1
	alice.x2 = alice.x2 ^ 1
	alice.x3 = alice.x3 ^ 1

	// Layer 2: And Bob's input with Alice's negated input
	d1_a := alice.x1 ^ alice.UVW[0].U
	d2_a := alice.x2 ^ alice.UVW[1].U
	d3_a := alice.x3 ^ alice.UVW[2].U

	e1_a := alice.y1 ^ alice.UVW[0].V
	e2_a := alice.y2 ^ alice.UVW[1].V
	e3_a := alice.y3 ^ alice.UVW[2].V

	d1_b := bob.x1 ^ bob.UVW[0].U
	d2_b := bob.x2 ^ bob.UVW[1].U
	d3_b := bob.x3 ^ bob.UVW[2].U

	e1_b := bob.y1 ^ bob.UVW[0].V
	e2_b := bob.y2 ^ bob.UVW[1].V
	e3_b := bob.y3 ^ bob.UVW[2].V

	// Send masked values to counterpart
	alice.Receive(d1_a, d2_a, d3_a)
	bob.Receive(e1_a, e2_a, e3_a)

	// Send masked values to counterpart
	b.Receive(d1, d2, d3)
	a.Receive(e1, e2, e3)

}
