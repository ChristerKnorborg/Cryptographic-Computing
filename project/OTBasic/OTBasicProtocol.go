package OTBasic

import (
	"crypto/rand"
	"cryptographic-computing/project/elgamal"
	"fmt"
	"log"
	mathRand "math/rand"
	"time"
)

// RandomBytes generates a slice of random bytes of a given length
func RandomBytes(length int) []byte {
	b := make([]byte, length)
	_, err := rand.Read(b)
	// Handle the error here. In production code, you might want to pass it up the call stack
	if err != nil {
		log.Fatal(err)
	}
	return b
}

// AllOnesBytes generates a slice of bytes of a given length, all set to 1
func AllOnesBytes(length int) []byte {
	b := make([]byte, length)
	for i := range b {
		b[i] = 1
	}
	return b
}

// AllTwosBytes generates a slice of bytes of a given length, all set to 2
func AllTwosBytes(length int) []byte {
	b := make([]byte, length)
	for i := range b {
		b[i] = 2
	}
	return b
}

// RandomSelectionBits generates a slice of m random selection bits (0 or 1)
func RandomSelectionBits(m int) []int {

	mathRand.NewSource(time.Now().UnixNano())

	bits := make([]int, m)
	for i := range bits {
		bits[i] = mathRand.Intn(2) // Generates a random integer 0 or 1
	}
	return bits
}

// k: Security parameter, l: Byte length of each message, m: Number of messages to be sent and selction bits
func OTBasicProtocol(l int, m int) {

	receiver := OTReceiver{}
	sender := OTSender{}

	elGamal := elgamal.ElGamal{}
	elGamal.Init() // Initialize the ElGamal encryption scheme

	selectionBits := RandomSelectionBits(m) // Receiver R holds m selection bits r = (r_1, ..., r_m).

	print("Selection bits: ")
	for _, b := range selectionBits {
		fmt.Printf("%d ", b) // Decimal print of []byte result
	}
	print("\n")

	// Initialize the sender with m pairs of messages
	// Generate m pairs of messages of l bytes each
	var messages []*MessagePair
	for i := 0; i < m; i++ {
		msg := MessagePair{
			Message0: RandomBytes(l),
			Message1: RandomBytes(l),
		}
		messages = append(messages, &msg)
	}
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
