package OTExtension

import (
	"cryptographic-computing/project/elgamal"
	"cryptographic-computing/project/utils"
	"fmt"
)

// k: Security parameter, l: Byte length of each message, m: Number of messages to be sent
func OTExtensionProtocol(k int, l int, m int, selectionBits []uint8, messages []*utils.MessagePair, elGamal elgamal.ElGamal) [][]byte {
	receiver := OTReceiver{}
	sender := OTSender{}

	// Initialize the receiver
	receiver.Init(selectionBits, k, l)

	sender.Init(messages, k, l)

	// Sender choose random string S. Receiver chooses seeds
	sender.ChooseRandomK()
	receiver.ChooseSeeds()

	// The parties invoke the κxOTκ_functionality (Sender plays receiver and receiver plays sender).
	// Original sender chooses a secret keys and public keys for each message, and sends public keys to original receiver, .
	// Original receiver chooses seeds and sends to original sender.

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

	if len(result) != len(messages) {
		fmt.Println("Result length is not equal to messages length in OTExtensionProtocol")
	}

	return result
}

// k: Security parameter, l: Byte length of each message, m: Number of messages to be sent
func OTExtensionProtocolEklundh(k int, l int, m int, selectionBits []uint8, messages []*utils.MessagePair, elGamal elgamal.ElGamal, multithreaded bool) [][]byte {
	receiver := OTReceiver{}
	sender := OTSender{}

	// Initialize the receiver
	receiver.Init(selectionBits, k, l)

	sender.Init(messages, k, l)

	// print("Messages: ")
	// for _, msg := range messages {
	// 	fmt.Printf("\n")
	// 	fmt.Printf("%d", msg.Message0)
	// 	fmt.Printf("%d", msg.Message1)
	// }
	// fmt.Printf("\n")

	// Sender choose random string S. Receiver chooses seeds
	sender.ChooseRandomK()
	receiver.ChooseSeeds()

	// The parties invoke the κxOTκ_functionality (Sender plays receiver and receiver plays sender).
	// Original sender chooses a secret keys and public keys for each message, and sends public keys to original receiver, .
	// Original receiver chooses seeds and sends to original sender.

	publicKeys := sender.Choose(&elGamal)
	receiver.ReceiveKeys(publicKeys)
	seedCiphertexts := receiver.EncryptSeeds(&elGamal)
	sender.DecryptSeeds(seedCiphertexts, &elGamal)

	// Receiver generates the Matrix T, and the Matrix U and send U to the sender.
	// The sender generates the Matrix Q from the received U Matrix.
	// Note, that every column q^i in Q is equal to (s_i * r) ⊕ t^i.
	U := receiver.GenerateMatrixTAndUEklundh(multithreaded)
	sender.GenerateMatrixQEklundh(U, multithreaded)

	// The sender sends (y0_j, y1_j) for every 1 ≤ j ≤ m to the receiver, where y0_j = x0_j ⊕ H(j, q_j) and y1_j = x1_j ⊕ H(j, q_j ⊕ s).
	// The receiver then computes x^(r_j)_j = y^(rj)_j ⊕ H(j, t_j) for every 1 ≤ j ≤ m. Then outputs (x^(r_1)_1, ..., x^(r_m)_m).
	ByteCiphertexts := sender.MakeAndSendCiphertexts()
	result := receiver.DecryptCiphertexts(ByteCiphertexts)

	if len(result) != len(messages) {
		fmt.Println("Result length is not equal to messages length in OTExtensionProtocol")
	}

	// print("Result: ")
	// for _, b := range result {
	// 	fmt.Printf("\n")
	// 	fmt.Printf("%d", b)
	// }

	return result

}

// k: Security parameter, l: Byte length of each message, m: Number of messages to be sent
func OTExtensionProtocolTranspose(k int, l int, m int, selectionBits []uint8, messages []*utils.MessagePair, elGamal elgamal.ElGamal) [][]byte {
	receiver := OTReceiver{}
	sender := OTSender{}

	// Initialize the receiver
	receiver.Init(selectionBits, k, l)

	sender.Init(messages, k, l)

	// Sender choose random string S. Receiver chooses seeds
	sender.ChooseRandomK()
	receiver.ChooseSeeds()

	// The parties invoke the κxOTκ_functionality (Sender plays receiver and receiver plays sender).
	// Original sender chooses a secret keys and public keys for each message, and sends public keys to original receiver, .
	// Original receiver chooses seeds and sends to original sender.

	publicKeys := sender.Choose(&elGamal)
	receiver.ReceiveKeys(publicKeys)
	seedCiphertexts := receiver.EncryptSeeds(&elGamal)
	sender.DecryptSeeds(seedCiphertexts, &elGamal)

	// Receiver generates the Matrix T, and the Matrix U and send U to the sender.
	// The sender generates the Matrix Q from the received U Matrix.
	// Note, that every column q^i in Q is equal to (s_i * r) ⊕ t^i.
	U := receiver.GenerateMatrixTAndUTranspose()
	sender.GenerateMatrixQTranspose(U)

	// The sender sends (y0_j, y1_j) for every 1 ≤ j ≤ m to the receiver, where y0_j = x0_j ⊕ H(j, q_j) and y1_j = x1_j ⊕ H(j, q_j ⊕ s).
	// The receiver then computes x^(r_j)_j = y^(rj)_j ⊕ H(j, t_j) for every 1 ≤ j ≤ m. Then outputs (x^(r_1)_1, ..., x^(r_m)_m).
	ByteCiphertexts := sender.MakeAndSendCiphertexts()
	result := receiver.DecryptCiphertexts(ByteCiphertexts)

	if len(result) != len(messages) {
		fmt.Println("Result length is not equal to messages length in OTExtensionProtocol")
	}

	return result

}
