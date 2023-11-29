package main

import (
	"cryptographic-computing/project/OTExtension"
	"cryptographic-computing/project/elgamal"
	utils "cryptographic-computing/project/utils"
	"math"
)

//Bench "cryptographic-computing/project/benchmark"

// func main() {
// 	//Bench.TestMakeDataFixL(24)
// 	utils.TestEklundhTranspose()
// 	//utils.TestDivideMatrix()
// }

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

	OTExtension.OTExtensionProtocolEklundh(k, l, m, selectionBits, messages, elGamal, false)

}
