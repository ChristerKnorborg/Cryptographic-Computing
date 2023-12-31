package handin4

func ObliviousTransfer(recipient Bloodtype, donor Bloodtype) bool {
	x := int(recipient) // Alice input (blood type)
	y := int(donor)     // Bob input (blood type)

	elGamal := ElGamal{}
	alice := Alice{}
	bob := Bob{}

	elGamal.InitFixedQ()     // initialize the ElGamal with public parameters p, q, g. Notice, q is fixed (Explained in the README)
	elGamal.TestProperties() // Test the properties of the ElGamal cryptosystem as described in Ivan Damgård's cryptography course book

	alice.Init(x) // Alice set her input x
	bob.Init(y)   // Bob set his input y

	publicKeys := alice.Choose(&elGamal)                  // Alice choose her input x and generate public keys - 7 fake and 1 real
	ciphertexts := bob.Transfer(publicKeys, &elGamal)     // Bob receives public keys, computes and transfers encrypted messages to Alice
	resultBigInt := alice.Retrieve(ciphertexts, &elGamal) // Alice receives the ciphertexts and decrypts the real message

	result := int(resultBigInt.Int64())
	return result == 1

}
