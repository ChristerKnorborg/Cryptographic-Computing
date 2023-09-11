package handin3

// One-Time Truth Table (OTTT) protocol that check if recipient blood type can receive donor blood type
func ComputeOTTTBloodTypeCompatability(recipient bloodtype, donor bloodtype) bool {
	AliceBloodType := int(recipient)
	BobBloodType := int(donor)
	
	n := 5 // number of and gates in our circuit ???

	d := Dealer{}
	a := Alice{}
	b := Bob{}
	
	d.Init(n) // Dealer generates two tuples of u, v and w for each AND gate - one for bob and one for alice

	a.Init(d.getAliceUVW()) // Alice gets the tuples of u's, v's and w's from the dealer
	b.Init(d.getBobUVW()) // Bob gets the tuples of u's, v's and w's from the dealer


	// Extract the bits from Alice (recipient)
	x1 := (recipient >> 2) & 1 // extract 3rd rightmost bit
	x2 := (recipient >> 1) & 1 // extract 2nd rightmost bit
	x3 := recipient & 1        // extract rightmost bit

	// Extract the bits from Bob (donor)
	y1 := (donor >> 2) & 1 // extract 3rd rightmost bit
	y2 := (donor >> 1) & 1 // extract 2nd rightmost bit
	y3 := donor & 1        // extract rightmost bit


	// Alice first layer
	






}



// BooleanFormula checks if recipient blood type can receive donor blood type using Boolean formulation
func BooleanFormula(recipient bloodtype, donor bloodtype) bool {

	x1 := (recipient >> 2) & 1 // extract 3rd rightmost bit
	x2 := (recipient >> 1) & 1 // extract 2nd rightmost bit
	x3 := recipient & 1        // extract rightmost bit

	y1 := (donor >> 2) & 1 // extract 3rd rightmost bit
	y2 := (donor >> 1) & 1 // extract 2nd rightmost bit
	y3 := donor & 1        // extract rightmost bit

	condition1 := (x1 == 0) || (y1 == 1)
	condition2 := (x2 == 0) || (y2 == 1)
	condition3 := (x3 == 0) || (y3 == 1)

	return condition1 && condition2 && condition3

}