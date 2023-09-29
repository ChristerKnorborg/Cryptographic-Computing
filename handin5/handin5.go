package handin5

func GarbledCircuit(recipient Bloodtype, donor Bloodtype) bool {
	x := int(recipient) // Alice input (blood type)
	y := int(donor)     // Bob input (blood type)

	alice := Alice{}
	bob := Bob{}
	elGamal := ElGamal{}

	input_write := 5

	elGamal.InitFixedQ() // initialize the ElGamal with public parameters p, q, g. Notice, q is fixed from handin4 (Explained in the README)
	alice.Init(x)        // Alice set her input x
	bob.Init(y)          // Bob set his input y

	// For each wire in the circuit, create 2 random strings of 128 bits

	// DUMMY
	return false
}
