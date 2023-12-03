package main

import "cryptographic-computing/project/benchmark"

// Outcomented benchmarking functions. Generates a csv file with the results.
func main() {
	//benchmark.TestMakeDataFixL(24)

	// //Make large matrix for Eklundh
	// k := int(math.Pow(2, float64(2))) // 2 ^ 7 = 128
	// M := int(math.Pow(2, float64(2))) // 2 ^ 26 = 67108864

	// matrix := make([][]byte, k)
	// for i := 0; i < k; i++ {
	// 	matrix[i] = utils.RandomSelectionBits(M)
	// }

	// benchmark.EklundhTranspose(matrix, false)

	k := 8 // 2^2 = 4 for a 4x4 matrix

	matrix := make([][]byte, k)
	for i := range matrix {
		matrix[i] = make([]byte, k)
		for j := range matrix[i] {
			matrix[i][j] = byte(i*k + j + 1)
		}
	}
	benchmark.EklundhTranspose(matrix, false)
}

// Main function can be used to test the different protocols manually
// and confirm that they output the correct messages depending on the selection bits.
// func main() {

// 	k := 128 // must be a power of 2
// 	l := 8
// 	m := int(math.Pow(2, float64(8)))

// 	// create cryptoalgorithm, messages and selection bits for algorithms.
// 	elGamal := elgamal.ElGamal{}
// 	elGamal.Init()
// 	selectionBits := utils.RandomSelectionBits(m)
// 	var messages []*utils.MessagePair
// 	for i := 0; i < m; i++ {
// 		msg := utils.MessagePair{
// 			Message0: utils.RandomBits(l),
// 			Message1: utils.RandomBits(l),
// 		}
// 		messages = append(messages, &msg)
// 	}

// 	/* Outcoment desired protocol that you want to test */

// 	//result := OTBasic.OTBasicProtocol(l, m, selectionBits, messages, elGamal)
// 	result := OTExtension.OTExtensionProtocol(k, l, m, selectionBits, messages, elGamal)
// 	//result := OTExtension.OTExtensionProtocolTranspose(k, l, m, selectionBits, messages, elGamal)
// 	//result := OTExtension.OTExtensionProtocolEklundh(k, l, m, selectionBits, messages, elGamal, false)
// 	//result := OTExtension.OTExtensionProtocolEklundh(k, l, m, selectionBits, messages, elGamal, true) // multithreaded

// 	println("")
// 	print("Selection bits: ")
// 	for _, b := range selectionBits {
// 		fmt.Printf("%d", b)
// 		print(" ")
// 	}
// 	fmt.Printf("\n")

// 	print("Messages: ")
// 	for _, msg := range messages {
// 		fmt.Printf("\n")
// 		fmt.Printf("%d", msg.Message0)
// 		fmt.Printf("%d", msg.Message1)
// 	}
// 	fmt.Printf("\n")

// 	print("Result: ")
// 	for _, b := range result {
// 		fmt.Printf("\n")
// 		fmt.Printf("%d", b)
// 	}
// }
