package main

import (
	"fmt"
	h5 "handin5"
	"math/big"
)

func main() {

	// // Try all combinations of blood types and compare with the lookup table
	// for recipient := h4.Ominus; recipient <= h4.ABplus; recipient++ {
	// 	for donor := h4.Ominus; donor <= h4.ABplus; donor++ {
	// 		aliceBloodType := recipient
	// 		bobBloodType := donor
	// 		ObliviousTransferResult := h4.GarbledCircuit(aliceBloodType, bobBloodType)
	// 		lookupTableResult := h4.LookUpBloodType(aliceBloodType, bobBloodType)
	// 		if ObliviousTransferResult != lookupTableResult {
	// 			fmt.Printf("Incorrect result for recipient: %s and donor: %s\n", h4.GetBloodTypeName(recipient), h4.GetBloodTypeName(donor))
	// 		} else {
	// 			fmt.Printf("Correct result for recipient: %s and donor: %s\n", h4.GetBloodTypeName(recipient), h4.GetBloodTypeName(donor))
	// 		}

	// 	}
	// }

	// Try a single combination of blood types
	aliceBloodType := h5.Ominus
	bobBloodType := h5.Aminus
	ObliviousTransferResult := h5.GarbledCircuit(aliceBloodType, bobBloodType)
	lookupTableResult := h5.LookUpBloodType(aliceBloodType, bobBloodType)
	if ObliviousTransferResult != lookupTableResult {
		fmt.Printf("Incorrect result for recipient: %s and donor: %s\n", h5.GetBloodTypeName(aliceBloodType), h5.GetBloodTypeName(bobBloodType))
	} else {
		fmt.Printf("Correct result for recipient: %s and donor: %s\n", h5.GetBloodTypeName(aliceBloodType), h5.GetBloodTypeName(bobBloodType))
	}

}

func testEncryptAndDecryptWithFixedValues() {

	x := 7

	str1 := "00000000000000000000000000000000"
	str2 := "00000000000000000000000000000000"
	testInput1 := h5.KeyPair{K_0: str1, K_1: str2}
	testInput2 := h5.KeyPair{K_0: str1, K_1: str2}
	testInput3 := h5.KeyPair{K_0: str1, K_1: str2}

	e_xInput := []h5.KeyPair{testInput1, testInput2, testInput3}

	elGamal := h5.ElGamal{}
	elGamal.InitFixedQ()

	secretKeys := make([]*big.Int, 3)
	for i := 0; i < 3; i++ {
		secretKeys[i] = elGamal.MakeSecretKey()
	}

	OTPublicKeys := h5.OTPublicKeys{}

	for i := 0; i < 3; i++ {
		OTPublicKeys.Keys[i][0] = elGamal.OGen()             // Fake key
		OTPublicKeys.Keys[i][1] = elGamal.Gen(secretKeys[i]) // Real key always 1 since we use 7 as input
	}

	encrypted_x := Encrypt(&elGamal, e_xInput, OTPublicKeys)
	result := Decrypt(&elGamal, secretKeys, encrypted_x, x)

	fmt.Printf(result[0])
	fmt.Println()
	fmt.Printf(result[1])
	fmt.Println()
	fmt.Printf(result[2])

}

func Encrypt(elGamal *h5.ElGamal, e_x []h5.KeyPair, OTKeys h5.OTPublicKeys) [3][2]*h5.Ciphertext {

	// Make 2 ciphertexts (of c1, c2) for each of Alice three input wires.
	// Two ciphertexts is due to encryption of both bits for each wire where one of the bits use is a real key and the other use a fake
	encrypted_x := [3][2]*h5.Ciphertext{}

	for i := 0; i < 3; i++ {

		keyString0, err0 := new(big.Int).SetString(e_x[i].K_0, 16) // Base 16 for hexadecimal string
		keyString1, err1 := new(big.Int).SetString(e_x[i].K_1, 16) // Base 16 for hexadecimal string

		if !err0 || !err1 {
			panic("Could not convert string to big.Int")
		}

		// Convert the strings to a big.Int
		wire_i_0 := keyString0 // Convert the first random string for Alice's input wire to big.Int
		wire_i_1 := keyString1 // Convert the second random string for Alice's input wire to big.Int

		encrypted_x[i][0] = elGamal.Encrypt(wire_i_0, OTKeys.Keys[i][0])
		encrypted_x[i][1] = elGamal.Encrypt(wire_i_1, OTKeys.Keys[i][1])

	}

	return encrypted_x

}

func Decrypt(elGamal *h5.ElGamal, sk []*big.Int, ciphertexts [3][2]*h5.Ciphertext, x int) []string {

	// Extract the bits from Alice's input
	inputInBits := h5.ExtractBits(x) // [x1, x2, x3]

	e_x := []string{}

	// Decrypt the ciphertexts
	for i := 0; i < 3; i++ {

		c1 := big.NewInt(0)
		c2 := big.NewInt(0)

		//
		if inputInBits[i] == 0 {
			c1 = ciphertexts[i][0].C1
			c2 = ciphertexts[i][0].C2

		} else if inputInBits[i] == 1 {
			c1 = ciphertexts[i][1].C1
			c2 = ciphertexts[i][1].C2
		}

		plaintextBigInt := elGamal.Decrypt(c1, c2, sk[i]) // Plaintext still in big int format
		plaintext := plaintextBigInt.Text(16)             // Plaintext in binary string format

		e_x = append(e_x, plaintext)
	}

	return e_x
}
