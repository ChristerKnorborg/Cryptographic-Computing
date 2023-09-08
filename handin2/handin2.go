package handin2

// One-Time Truth Table (OTTT) protocol that check if recipient blood type can receive donor blood type
func ComputeOTTTBloodTypeCompatability(recipient bloodtype, donor bloodtype) bool {
	AliceBloodType := recipient
	BobBloodType := donor

	d := Dealer{}
	d.Init() // Dealer initializes the matrices and coordinates

	a := Alice{}
	a.Init(int(AliceBloodType), d.GetMatrixA(), d.GetR()) // Alice initializes her matrix and coordinates

	b := Bob{}
	b.Init(int(BobBloodType), d.GetMatrixB(), d.GetS()) // Bob initializes his matrix and coordinates

	b.Receive(a.Send()) // Immitate Alice sends u to Bob

	v, z_B := b.Send() // Immitate Bob sends v and z_B to Alice
	a.Receive(v, z_B)

	return a.ComputeOutput() // Alice computes output z = M_A[u, v] ⊕ z_B (same as f(x,y)) and returns it
}
