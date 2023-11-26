package OTBasic

import (
	"cryptographic-computing/project/elgamal"
	"cryptographic-computing/project/utils"
	"fmt"
)

// k: Security parameter, l: Byte length of each message, m: Number of messages to be sent and selction bits
func OTBasicProtocol(l int, m int, selectionBits []uint8, messages []*utils.MessagePair, elGamal elgamal.ElGamal) {

	receiver := OTReceiver{}
	sender := OTSender{}

	// Initialize the receiver's selection bits and the sender's messages
	receiver.Init(selectionBits)
	sender.Init(messages, m)

	fmt.Println("SelectionBits Basic: ")
	for _, b := range selectionBits {
		fmt.Printf("%d ", b)
	}
	fmt.Println()

	fmt.Println("Messages Basic: ")
	for _, message := range messages {
		for _, b := range message.Message0 {
			fmt.Printf("%d ", b)
		}
		fmt.Print(" ")

		for _, b := range message.Message1 {
			fmt.Printf("%d ", b)
		}
		fmt.Println()
	}

	// The receiver makes secret keys and oblivious keys each message to be received based on the selection bits.
	// Then send the public keys to the sender.
	publicKeys := receiver.Choose(m, &elGamal)
	sender.ReceiveKeys(publicKeys)

	// Sender encrypts the messages using the public keys received from the receiver.
	// Then send the ciphertexts to the receiver.
	ciphertextPairs := sender.EncryptMessages(&elGamal)

	// The receiver decrypts the ciphertexts using the secret keys depending on the selection bits.
	result := receiver.DecryptMessage(ciphertextPairs, l, &elGamal)

	if len(result) != len(messages) {
		fmt.Println("Result length is not equal to messages length in OTBasicProtocol")
	}

	print("Result Basic: ")
	for _, b := range result {
		fmt.Printf("%d ", b)
	}
	fmt.Println()
}
