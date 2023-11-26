package main

import (
	OTExt "cryptographic-computing/project/OTExtension"
	"cryptographic-computing/project/elgamal"
	"cryptographic-computing/project/utils"
	"math"
)

// func main() {
// 	//Bench.TestMakeDataFixL(3)
// 	//Util.TestEklundhTranspose()
// }

func main() {

	k := 256
	l := 1
	m := int(math.Pow(2, float64(1)))

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

	OTExt.OTExtensionProtocol(k, l, m, selectionBits, messages, elGamal)
}
