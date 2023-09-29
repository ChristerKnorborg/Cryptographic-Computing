package handin5

func GarbledCircuit(recipient Bloodtype, donor Bloodtype) {
	x := int(recipient) // Alice input (blood type)
	y := int(donor)     // Bob input (blood type)

	alice := Alice{}
	bob := Bob{}
	elGamal := ElGamal{}

	elGamal.InitFixedQ() // initialize the ElGamal with public parameters p, q, g. Notice, q is fixed from handin4 (Explained in the README)
	alice.Init(x)        // Alice set her input x
	bob.Init(y)          // Bob set his input y

}
