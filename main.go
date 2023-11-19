package main

import (
	"cryptographic-computing/project/OTExtension"
	"cryptographic-computing/project/elgamal"
)

func main() {
	receiver := OTExtension.OTReceiver{}
	sender := OTExtension.OTSender{}

	// Initialize the receiver
	selectionBits := []int{0, 1, 0, 1, 0, 1} // Receiver R holds m selection bits r = (r_1, ..., r_m).
	receiver.Init(selectionBits)

	// Initialize the sender with m pairs of messages
	msg1 := OTExtension.MessagePair{Message0: "alice", Message1: "bob"}
	msg2 := OTExtension.MessagePair{Message0: "alice", Message1: "bob"}
	msg3 := OTExtension.MessagePair{Message0: "alice", Message1: "bob"}
	msg4 := OTExtension.MessagePair{Message0: "alice", Message1: "bob"}
	msg5 := OTExtension.MessagePair{Message0: "alice", Message1: "bob"}
	msg6 := OTExtension.MessagePair{Message0: "alice", Message1: "bob"}
	messages := []*OTExtension.MessagePair{&msg1, &msg2, &msg3, &msg4, &msg5, &msg6}
	selectionBits_κxOTκ_functionality := []int{0, 1}
	k := 2

	sender.Init(messages, k, selectionBits_κxOTκ_functionality)

	// Sender choose random string S. Receiver chooses seeds
	sender.ChooseRandomString()
	receiver.ChooseSeeds(k)

	// The parties invoke the κxOTκ_functionality (Sender plays receiver and receiver plays sender).
	// Original sender chooses public keys and sends to original receiver.
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
	receiver.GenerateMatrixT()
	U := receiver.GenerateAndSendUMatrix()
	sender.GenerateQMatrix(U)

}
