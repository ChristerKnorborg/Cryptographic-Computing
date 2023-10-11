package handin6

func HomomorphicBloodtypeEncryption(recipient Bloodtype, donor Bloodtype) bool {
	x := int(recipient) // Alice input (blood type)
	y := int(donor)     // Bob input (blood type)

	alice := Alice{}
	bob := Bob{}
	DHE := DHE{} // Depth d Homomorphic Encryption (DHE) scheme

	alice.Init(x) // Alice set her input x to be the bits x1, x2, x3
	bob.Init(y)   // Bob set his input y to be the bits y1, y2, y3

	bitlenP := 512                             // bit length of p
	DHE.GenerateKeys(bitlenP)                  // Generate the parameters p, q = (q_1,..., q_n), r = (r_1,..., r_n), y = (y_1,..., y_n) for the DHE scheme
	alice.RecieveSecretKey(DHE.GetSecretKey()) // Alice recieves the secret key p

	encryptedX := alice.Encrypt(&DHE) // Alice encrypts her input bits x1, x2, x3 using DHE and sends them to Bob
	bob.Encrypt(&DHE)                 // Bob encrypts his input bits y1, y2, y3 using DHE (Stores locally for bob)

	evaluatedOutput := bob.RecieveAndEvaluate(encryptedX) // Bob recieves Alice's encrypted input bits and evaluates the circuit

	result := alice.Decrypt(evaluatedOutput, &DHE) // Alice decrypts the result of the circuit evaluation using DHE with her secret key p

	return result == 1 // 1 = true, 0 = false

}
