package OTBasic

import (
	"cryptographic-computing/project/elgamal"
	"fmt"
)

// k: Security parameter, l: Byte length of each message, m: Number of messages to be sent and selction bits
func OTBasicProtocol(l int, m int, selectionBits []int, messages []*MessagePair, elGamal elgamal.ElGamal) {

	receiver := OTReceiver{}
	sender := OTSender{}

	print("Selection bits: ")
	for _, b := range selectionBits {
		fmt.Printf("%d ", b) // Decimal print of []byte result
	}
	print("\n")

	// Initialize the sender with m pairs of messages
	// Generate m pairs of messages of l bytes each
	print("Messages: ")
	for _, b := range messages {
		fmt.Printf("%d ", b.Message0) // Decimal print of []byte result
		fmt.Printf("%d ", b.Message1) // Decimal print of []byte result
	}
	print("\n")

	// Initialize the receiver's selection bits and the sender's messages
	receiver.Init(selectionBits)
	sender.Init(messages, m)

	// The receiver makes secret keys and oblivious keys each message to be received based on the selection bits.
	// Then send the public keys to the sender.
	publicKeys := receiver.Choose(m, &elGamal)
	sender.ReceiveKeys(publicKeys)

	// Sender encrypts the messages using the public keys received from the receiver.
	// Then send the ciphertexts to the receiver.
	ciphertextPairs := sender.EncryptMessages(&elGamal)

	// The receiver decrypts the ciphertexts using the secret keys depending on the selection bits.
	result := receiver.DecryptMessage(ciphertextPairs, l, &elGamal)

	print("Result: ")
	for _, b := range result {
		fmt.Printf("%d ", b) // Decimal print of []byte result
	}

}
