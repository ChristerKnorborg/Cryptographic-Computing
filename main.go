package main

import (
	"cryptographic-computing/project/OTExtension"
	"cryptographic-computing/project/elgamal"
	"fmt"
)

func main() {
	receiver := OTExtension.OTReceiver{}
	sender := OTExtension.OTSender{}
	k := 2 // Security parameter
	l := 1 // byte lenght of each message

	// Initialize the receiver
	selectionBits := []int{0, 1, 0, 1, 0, 1} // Receiver R holds m selection bits r = (r_1, ..., r_m).
	receiver.Init(selectionBits, k, l)

	// Initialize the sender with m pairs of messages
	msg1 := OTExtension.MessagePair{Message0: []byte{1}, Message1: []byte{0}}
	msg2 := OTExtension.MessagePair{Message0: []byte{1}, Message1: []byte{0}}
	msg3 := OTExtension.MessagePair{Message0: []byte{1}, Message1: []byte{0}}
	msg4 := OTExtension.MessagePair{Message0: []byte{1}, Message1: []byte{0}}
	msg5 := OTExtension.MessagePair{Message0: []byte{1}, Message1: []byte{0}}
	msg6 := OTExtension.MessagePair{Message0: []byte{1}, Message1: []byte{0}}
	messages := []*OTExtension.MessagePair{&msg1, &msg2, &msg3, &msg4, &msg5, &msg6}

	// DEBUG
	// for _, m := range messages {
	// 	fmt.Printf("%d ", m.Message0)
	// 	fmt.Printf("%d ", m.Message1)
	// 	fmt.Println()
	// }

	sender.Init(messages, k, l)

	// Sender choose random string S. Receiver chooses seeds
	sender.ChooseRandomString()
	receiver.ChooseSeeds()

	// The parties invoke the κxOTκ_functionality (Sender plays receiver and receiver plays sender).
	// Original sender chooses a secret keys and public keys for each message, and sends public keys to original receiver, .
	// Original receiver chooses seeds and sends to original sender.

	elGamal := elgamal.ElGamal{}
	elGamal.Init()

	publicKeys := sender.Choose(&elGamal)
	receiver.ReceiveKeys(publicKeys)
	seedCiphertexts := receiver.EncryptSeeds(&elGamal)
	sender.DecryptSeeds(seedCiphertexts, &elGamal)

	// Receiver generates the Matrix T, and the Matrix U and send U to the sender.
	// The sender generates the Matrix Q from the received U Matrix.
	// Note, that every column q^i in Q is equal to (s_i * r) ⊕ t^i.
	receiver.GenerateMatrixT()
	U := receiver.GenerateAndSendMatrixU()
	sender.GenerateMatrixQ(U)

	// The sender sends (y0_j, y1_j) for every 1 ≤ j ≤ m to the receiver, where y0_j = x0_j ⊕ H(j, q_j) and y1_j = x1_j ⊕ H(j, q_j ⊕ s).
	// The receiver then computes x^(r_j)_j = y^(rj)_j ⊕ H(j, t_j) for every 1 ≤ j ≤ m. Then outputs (x^(r_1)_1, ..., x^(r_m)_m).
	ByteCiphertexts := sender.MakeAndSendCiphertexts()
	result := receiver.DecryptCiphertexts(ByteCiphertexts)

	for _, b := range result {
		fmt.Printf("%d ", b) // Decimal print of []byte result
	}
}
