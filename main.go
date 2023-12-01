package main

import (
	"cryptographic-computing/project/OTExtension"
	"cryptographic-computing/project/elgamal"
	"cryptographic-computing/project/utils"
	"fmt"
	"math"
)

// Outcomented benchmarking functions. Generates a csv file with the results.
// func main() {
// 	//Bench.TestMakeDataFixL(24)
// }

// Main function can be used to test the different protocols manually
// and confirm that they output the correct messages depending on the selection bits.
func main() {

	k := 2 // must be a power of 2
	l := 1
	m := int(math.Pow(2, float64(2)))

	// create cryptoalgorithm, messages and selection bits for algorithms.
	elGamal := elgamal.ElGamal{}
	elGamal.Init()
	selectionBits := utils.RandomSelectionBits(m)
	var messages []*utils.MessagePair
	for i := 0; i < m; i++ {
		msg := utils.MessagePair{
			Message0: utils.RandomBytes(l),
			Message1: utils.RandomBytes(l),
		}
		messages = append(messages, &msg)
	}

	/* Outcoment desired protocol that you want to test */

	//OTBasicResult := OTBasic.OTBasicProtocol(l, m, selectionBits, messages, elGamal)
	//OTExtensionResult := OTExtension.OTExtensionProtocol(k, l, m, selectionBits, messages, elGamal)
	//OTExtensionTransposeResult := OTExtension.OTExtensionProtocolTranspose(k, l, m, selectionBits, messages, elGamal)
	OTExtensionEklundhResult := OTExtension.OTExtensionProtocolEklundh(k, l, m, selectionBits, messages, elGamal, false)
	//OTExtensionEklundhResultMultitreaded := OTExtension.OTExtensionProtocolEklundh(k, l, m, selectionBits, messages, elGamal, true)

	println("")
	print("Selection bits: ")
	for _, b := range selectionBits {
		fmt.Printf("%d", b)
		print(" ")
	}
	fmt.Printf("\n")

	print("Messages: ")
	for _, msg := range messages {
		fmt.Printf("\n")
		fmt.Printf("%d", msg.Message0)
		fmt.Printf("%d", msg.Message1)
	}
	fmt.Printf("\n")

	print("Result: ")
	for _, b := range OTExtensionEklundhResult {
		fmt.Printf("\n")
		fmt.Printf("%d", b)
	}
}
