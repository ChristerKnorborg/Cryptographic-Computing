package OTExtension

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

// k: Security parameter, l: Byte length of each message, m: Number of messages to be sent
func OTExtensionProtocol(k int, l int, m int) {
	receiver := OTReceiver{}
	sender := OTSender{}

	// Initialize the receiver
	selectionBits := RandomSelectionBits(m) // Receiver R holds m selection bits r = (r_1, ..., r_m).
	receiver.Init(selectionBits, k, l)

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

	sender.Init(messages, k, l)

	// Sender choose random string S. Receiver chooses seeds
	sender.ChooseRandomK()
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
