package main

import (
	"cryptographic-computing/project/OTExtension"
	"cryptographic-computing/project/elgamal"
)

func main() {
	receiver := OTExtension.OTReceiver{}
	sender := OTExtension.OTSender{}
	k := 2 // Security parameter

	// Initialize the receiver
	selectionBits := []int{0, 1, 0, 1, 0, 1} // Receiver R holds m selection bits r = (r_1, ..., r_m).
	receiver.Init(selectionBits, k)

	// Initialize the sender with m pairs of messages
	msg1 := OTExtension.MessagePair{Message0: []byte{0}, Message1: []byte{1}}
	msg2 := OTExtension.MessagePair{Message0: []byte{0}, Message1: []byte{1}}
	msg3 := OTExtension.MessagePair{Message0: []byte{0}, Message1: []byte{1}}
	msg4 := OTExtension.MessagePair{Message0: []byte{0}, Message1: []byte{1}}
	msg5 := OTExtension.MessagePair{Message0: []byte{0}, Message1: []byte{1}}
	msg6 := OTExtension.MessagePair{Message0: []byte{0}, Message1: []byte{1}}
	messages := []*OTExtension.MessagePair{&msg1, &msg2, &msg3, &msg4, &msg5, &msg6}
	selectionBits_κxOTκ_functionality := []int{0, 1}

	sender.Init(messages, k, selectionBits_κxOTκ_functionality)

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
	ciphertexts := receiver.EncryptSeeds(&elGamal)
	sender.DecryptSeeds(ciphertexts, &elGamal)

	// Receiver generates the Matrix T, and the Matrix U and send U to the sender.
	// The sender generates the Matrix Q from the received U Matrix.
	// Note, that every column q^i in Q is equal to (s_i * r) ⊕ t^i.

	// receiver.GenerateMatrixT()
	// U := receiver.GenerateAndSendUMatrix()
	// sender.GenerateQMatrix(U)

}
