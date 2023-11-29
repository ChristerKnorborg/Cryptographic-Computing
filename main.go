package main

import utils "cryptographic-computing/project/utils"

//Bench "cryptographic-computing/project/benchmark"

func main() {
	//Bench.TestMakeDataFixL(24)
	//utils.TestEklundhTranspose()
	utils.TestEklundhTranspose()
}

// func main() {

// 	k := 128 // must be a power of 2
// 	l := 1
// 	m := int(math.Pow(2, float64(16)))

// 	// create cryptoalgorithm, messages and selection bits for algorithms.
// 	elGamal := elgamal.ElGamal{}
// 	elGamal.Init()
// 	selectionBits := utils.RandomSelectionBits(m)
// 	var messages []*utils.MessagePair
// 	for i := 0; i < m; i++ {
// 		msg := utils.MessagePair{
// 			Message0: utils.RandomBytes(l),
// 			Message1: utils.RandomBytes(l),
// 		}
// 		messages = append(messages, &msg)
// 	}

// 	OTExtension.OTExtensionProtocolEklundh(k, l, m, selectionBits, messages, elGamal)

// }
