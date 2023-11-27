package main

import (
	Bench "cryptographic-computing/project/benchmark"
	//Util "cryptographic-computing/project/utils"
)

func main() {
	Bench.TestMakeDataFixL(16)
	//Util.TestEklundhTranspose()
}

// func main() {

// 	k := 128
// 	l := 1
// 	m := int(math.Pow(2, float64(8)))

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

// 	OTExt.OTExtensionProtocolEklundh(k, l, m, selectionBits, messages, elGamal)
// }
