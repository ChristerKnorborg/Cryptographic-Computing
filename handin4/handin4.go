package handin4

func ObliviousTransfer(recipient Bloodtype, donor Bloodtype) bool {
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

	elGamal := ElGamal{}
	alice := Alice{}
	bob := Bob{}

	elGamal.Init() // initialize the ElGamal with public parameters p, q, g

	alice.Init(elGamal.p, elGamal.q, elGamal.g) // initialize Alice with public parameters p, q, g
	bob.Init(elGamal.p, elGamal.q, elGamal.g)   // initialize Bob with public parameters p, q, g

	alice.MakePublicKeys(x1, x2, x3) // Alice inputs her bits
	bob.Input(y1, y2, y3)            // Bob inputs his bits
}
