package handin6

func HomomorphicBloodtypeEncryption(recipient Bloodtype, donor Bloodtype) bool {
	x := int(recipient) // Alice input (blood type)
	y := int(donor)     // Bob input (blood type)

	alice := Alice{}
	bob := Bob{}
	DHE := DHE{} // Depth d Homomorphic Encryption (DHE) scheme

	DHE.GenerateKeys() // Generate the parameters p, q = (q_1,..., q_n), r = (r_1,..., r_n), y = (y_1,..., y_n) for the DHE scheme

	encryptedX := alice.Init(x, &DHE) // Alice set her input x
	bob.Init(y, &DHE)                 // Bob set his input y

	evaluatedOutput := bob.RecieveAndEvaluate(encryptedX)

	result := alice.Decide(evaluatedOutput, &DHE)

	return result == 1 // 1 = true, 0 = false

}
