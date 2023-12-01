package OTExtension

import (
	"cryptographic-computing/project/elgamal"
	"cryptographic-computing/project/utils"
	"fmt"
)

func OTExtensionProtocol(k int, l int, m int, selectionBits []byte, messages []*utils.MessagePair, elGamal elgamal.ElGamal) [][]byte {
	receiver := OTReceiver{}
	sender := OTSender{}

	// Initialize public parameters for both parties, the receiver's selection bits, and the sender's messages
	receiver.Init(selectionBits, k, l)
	sender.Init(messages, k, l)

	// Sender choose random string S. Receiver chooses k random seeds. All of length k.
	sender.ChooseRandomS()
	receiver.ChooseSeeds()

	// The parties invoke the regular OT functionality k times (Sender plays receiver and receiver plays sender).
	// Sender chooses a secret keys and public keys for each message, and sends public keys to original receiver, .
	// Receiver chooses seeds and sends to original sender.
	publicKeys := sender.Choose(&elGamal)
	receiver.ReceiveKeys(publicKeys)
	seedCiphertexts := receiver.EncryptSeeds(&elGamal)
	sender.DecryptSeeds(seedCiphertexts, &elGamal)

	// Receiver generates the Matrix T, and the Matrix U and send U to the sender.
	// The sender generates the Matrix Q from the received U Matrix.
	U := receiver.GenerateMatrixTAndU()
	sender.GenerateMatrixQ(U)

	// The sender sends m ciphertext pairs to the receiver.
	// The receiver computes the desired message based on the selection bits.
	ByteCiphertexts := sender.MakeAndSendCiphertexts()
	result := receiver.DecryptCiphertexts(ByteCiphertexts)

	if len(result) != len(messages) {
		fmt.Println("Result length is not equal to messages length in OTExtensionProtocol")
	}

	return result
}

func OTExtensionProtocolTranspose(k int, l int, m int, selectionBits []byte, messages []*utils.MessagePair, elGamal elgamal.ElGamal) [][]byte {
	receiver := OTReceiver{}
	sender := OTSender{}

	// Initialize public parameters for both parties, the receiver's selection bits, and the sender's messages
	receiver.Init(selectionBits, k, l)
	sender.Init(messages, k, l)

	// Sender choose random string S. Receiver chooses k random seeds. All of length k.
	sender.ChooseRandomS()
	receiver.ChooseSeeds()

	// The parties invoke the regular OT functionality k times (Sender plays receiver and receiver plays sender).
	// Sender chooses a secret keys and public keys for each message, and sends public keys to original receiver, .
	// Receiver chooses seeds and sends to original sender.
	publicKeys := sender.Choose(&elGamal)
	receiver.ReceiveKeys(publicKeys)
	seedCiphertexts := receiver.EncryptSeeds(&elGamal)
	sender.DecryptSeeds(seedCiphertexts, &elGamal)

	// Receiver generates the Matrix T, and the Matrix U and send U to the sender.
	// The sender generates the Matrix Q from the received U Matrix.
	U := receiver.GenerateMatrixTAndUTranspose()
	sender.GenerateMatrixQTranspose(U)

	// The sender sends m ciphertext pairs to the receiver.
	// The receiver computes the desired message based on the selection bits.
	ByteCiphertexts := sender.MakeAndSendCiphertexts()
	result := receiver.DecryptCiphertexts(ByteCiphertexts)

	if len(result) != len(messages) {
		fmt.Println("Result length is not equal to messages length in OTExtensionProtocol")
	}

	return result
}

func OTExtensionProtocolEklundh(k int, l int, m int, selectionBits []byte, messages []*utils.MessagePair, elGamal elgamal.ElGamal, multithreaded bool) [][]byte {
	receiver := OTReceiver{}
	sender := OTSender{}

	// Initialize public parameters for both parties, the receiver's selection bits, and the sender's messages
	receiver.Init(selectionBits, k, l)
	sender.Init(messages, k, l)

	// Sender choose random string S. Receiver chooses k random seeds. All of length k.
	sender.ChooseRandomS()
	receiver.ChooseSeeds()

	// The parties invoke the regular OT functionality k times (Sender plays receiver and receiver plays sender).
	// Sender chooses a secret keys and public keys for each message, and sends public keys to original receiver, .
	// Receiver chooses seeds and sends to original sender.
	publicKeys := sender.Choose(&elGamal)
	receiver.ReceiveKeys(publicKeys)
	seedCiphertexts := receiver.EncryptSeeds(&elGamal)
	sender.DecryptSeeds(seedCiphertexts, &elGamal)

	// Receiver generates the Matrix T, and the Matrix U and send U to the sender.
	// The sender generates the Matrix Q from the received U Matrix.
	U := receiver.GenerateMatrixTAndUEklundh(false)
	sender.GenerateMatrixQEklundh(U, false)

	// The sender sends m ciphertext pairs to the receiver.
	// The receiver computes the desired message based on the selection bits.
	ByteCiphertexts := sender.MakeAndSendCiphertexts()
	result := receiver.DecryptCiphertexts(ByteCiphertexts)

	if len(result) != len(messages) {
		fmt.Println("Result length is not equal to messages length in OTExtensionProtocol")
	}

	return result

}
