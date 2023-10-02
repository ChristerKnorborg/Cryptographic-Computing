package handin5

func GarbledCircuit(recipient Bloodtype, donor Bloodtype) bool {
	x := int(recipient) // Alice input (blood type)
	y := int(donor)     // Bob input (blood type)

	alice := Alice{}
	bob := Bob{}
	elGamal := ElGamal{}

	elGamal.InitFixedQ() // initialize the ElGamal with public parameters p, q, g. Notice, q is fixed from handin4 (Explained in the README for handin4)
	alice.Init(x)        // Alice set her input x
	bob.Init(y)          // Bob set his input y

	// Bob runs Gb to make the garbled circuit [F]. d is the Z values from the output of the garbled circuit,
	// e_xor is the XOR values from the output of the garbled circuit
	F, d, e_xor := bob.MakeGarbledCircuit()
	Y := bob.Encode() // Bob runs En to encode his input y to make [Y]

	/* OBLIVIOUS TRANSFER */
	// Alice and Bob run a secure two party computation, where Bob inputs e_x, Alice inputs x, and Alice learns [X]
	publicKeys := alice.MakeAndTransferKeys(&elGamal) // Alice uses ElGamal to make herself a secret key, and two public keys for Bob
	bob.ReceiveKeys(publicKeys)                       // Bob receives the public keys from Alice
	e_xCiphertext := bob.Encrypt(&elGamal)            // Bob encrypts his ?wire values? with the public keys from Alice
	alice.Decrypt(e_xCiphertext, &elGamal)            // Alice decrypts the ?wire values? from Bob and stores them locally

	// Alice receives the encrypted values from Bob and evaluates the garbled circuit to get result Z′.
	// She outputs z = 0 if Z′ = Z_0, z = 1 if Z′ = Z_1 or Panic (due to non-honest opposition) if Z′ ∉ {Z0, Z1}.
	z := alice.EvaluateGarbledCircuit(F, d, e_xor, Y, bob.e_xor)

	return z == 1 // returns true if Alice can receive blood from Bob, otherwise false
}
