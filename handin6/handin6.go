package handin6

func GarbledCircuit(recipient Bloodtype, donor Bloodtype) bool {
	x := int(recipient) // Alice input (blood type)
	y := int(donor)     // Bob input (blood type)

	alice := Alice{}
	bob := Bob{}
	DHE := DHE{} // Depth d Homomorphic Encryption (DHE) scheme

	DHE.GenerateKeys() // Generate the parameters p, q = (q_1,..., q_n), r = (r_1,..., r_n), y = (y_1,..., y_n) for the DHE scheme
	DHE

	alice.Init(x) // Alice set her input x
	bob.Init(y)   // Bob set his input y

	// Bob runs Gb to make the garbled circuit [F]. d is the Z values from the output of the garbled circuit,
	// e_xor is the XOR wires in F. Bob also runs En to encode his input wires y to make [Y]
	F, d, e_xor := bob.MakeGarbledCircuit()
	Y := bob.Encode()

	/* OBLIVIOUS TRANSFER */
	// Alice and Bob run a secure two party computation, where Bob inputs e_x, Alice inputs x, and Alice learns [X].
	// Alice uses ElGamal to make herself a secret key, and two public keys (1 real and 1 fake) for Bob per input bit.
	// That is she does 1 over 2 OT for each of her three input bits.
	publicKeys := alice.MakeAndTransferKeys(&elGamal)
	bob.ReceiveKeys(publicKeys)  // Bob receives the public keys from Alice
	e_x := bob.Encrypt(&elGamal) // Bob encrypts the wires corresponding to Alice's input bits with the public keys from Alice
	alice.Decrypt(e_x, &elGamal) // Alice decrypts the wires corresponding to Alice's input bits and stores them locally

	// Alice receives the encrypted values from Bob and evaluates the garbled circuit to get result Z′.
	// She outputs z = 0 if Z′ = Z_0, z = 1 if Z′ = Z_1 or Panic (due to non-honest opposition) if Z′ ∉ {Z0, Z1}.
	z := alice.EvaluateGarbledCircuit(F, d, e_xor, Y)

	return z == 1 // returns true if Alice can receive blood from Bob, otherwise false
}
