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

	a.Init(d.getAliceUVW()) // Alice gets the tuples of u, v and w from the dealer
	b.Init(d.getBobUVW()) // Bob gets the tuples of u, v and w from the dealer

	

}
