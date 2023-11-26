package OTExtension

import (
	"cryptographic-computing/project/elgamal"
	"cryptographic-computing/project/utils"
	"fmt"
	"strconv"
)

// k: Security parameter, l: Byte length of each message, m: Number of messages to be sent
func OTExtensionProtocol(k int, l int, m int, selectionBits []uint8, messages []*utils.MessagePair, elGamal elgamal.ElGamal) {
	receiver := OTReceiver{}
	sender := OTSender{}

	// Initialize the receiver
	receiver.Init(selectionBits, k, l)

	sender.Init(messages, k, l)

	fmt.Println("SelectionBits Extension: ")
	for _, b := range selectionBits {
		fmt.Printf("%d ", b)
	}
	fmt.Println()

	fmt.Println("Messages: ")
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
	T := receiver.GenerateMatrixT()
	U := receiver.GenerateAndSendMatrixU()
	Q := sender.GenerateMatrixQ(U)

	print("CHECKING Q MATRIX COLS: \n")
	for i := 0; i < k; i++ {
		s_i := sender.S[i]
		for j := 0; j < m; j++ {
			q_idx := Q[j][i]
			r_j := receiver.SelectionBits[j]
			t_i := T[j][i]
			if q_idx != (s_i*r_j)^t_i {
				G := q_idx ^ (s_i * r_j)
				print("s_i, u[j][i], G(): " + strconv.Itoa(int(s_i)) + " " + strconv.Itoa(int(U[j][i])) + " " + strconv.Itoa(int(G)) + "\n")
				print("s_i, r_j, t_i: " + strconv.Itoa(int(s_i)) + " " + strconv.Itoa(int(r_j)) + " " + strconv.Itoa(int(t_i)) + "\n")
				print("q_idx: " + strconv.Itoa(int(q_idx)) + "\n")
				fmt.Println("ERROR: q_idx != (s_i * r_j) ^ t_i")

			}
		}
	}

	// DEBUG: CHECK THAT EVERY ROW IN Q IS EQUAL TO (r_j*s) ⊕ t_j
	print("CHECKING Q MATRIX ROWS: \n")
	for j := 0; j < m; j++ {
		for i := 0; i < k; i++ {
			q_idx := Q[j][i]
			r_j := receiver.SelectionBits[j]
			s_i := sender.S[i]
			t_j := T[j][i]

			if q_idx != (r_j*s_i)^t_j {
				print("r_j, s_i, t[j][i]: " + strconv.Itoa(int(r_j)) + " " + strconv.Itoa(int(s_i)) + " " + strconv.Itoa(int(t_j)) + "\n")
				print("q_idx: " + strconv.Itoa(int(q_idx)) + "\n")
				fmt.Println("ERROR: q_idx != (r_j * s_i) ^ t_j")
			}
		}
	}

	// The sender sends (y0_j, y1_j) for every 1 ≤ j ≤ m to the receiver, where y0_j = x0_j ⊕ H(j, q_j) and y1_j = x1_j ⊕ H(j, q_j ⊕ s).
	// The receiver then computes x^(r_j)_j = y^(rj)_j ⊕ H(j, t_j) for every 1 ≤ j ≤ m. Then outputs (x^(r_1)_1, ..., x^(r_m)_m).
	ByteCiphertexts := sender.MakeAndSendCiphertexts()
	result := receiver.DecryptCiphertexts(ByteCiphertexts)

	if len(result) != len(messages) {
		fmt.Println("Result length is not equal to messages length in OTExtensionProtocol")
	}

	print("Result: ")
	for _, b := range result {
		fmt.Printf("%d ", b)
	}
	fmt.Println()
}

// k: Security parameter, l: Byte length of each message, m: Number of messages to be sent
func OTExtensionProtocolEklundh(k int, l int, m int, selectionBits []uint8, messages []*utils.MessagePair, elGamal elgamal.ElGamal) {
	receiver := OTReceiver{}
	sender := OTSender{}

	// Initialize the receiver
	receiver.Init(selectionBits, k, l)

	sender.Init(messages, k, l)

	// Print messages as decimal numbers
	fmt.Println("Messages: ")
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
	U, T := receiver.GenerateMatrixTAndUEklundh()
	Q := sender.GenerateMatrixQEklundh(U)

	U = utils.EklundhTransposeMatrix(U) // Debug (We never transpose U in the protocol as it is not needed)
	print("CHECKING Q MATRIX COLS: \n")
	for i := 0; i < k; i++ {
		s_i := sender.S[i]
		for j := 0; j < m; j++ {
			q_idx := Q[j][i]
			r_j := receiver.SelectionBits[j]
			t_i := T[j][i]
			if q_idx != (s_i*r_j)^t_i {
				G := q_idx ^ (s_i * r_j)
				print("s_i, u[j][i], G(): " + strconv.Itoa(int(s_i)) + " " + strconv.Itoa(int(U[j][i])) + " " + strconv.Itoa(int(G)) + "\n")
				print("s_i, r_j, t_i: " + strconv.Itoa(int(s_i)) + " " + strconv.Itoa(int(r_j)) + " " + strconv.Itoa(int(t_i)) + "\n")
				print("q_idx: " + strconv.Itoa(int(q_idx)) + "\n")
				fmt.Println("ERROR: q_idx != (s_i * r_j) ^ t_i")

			}
		}
	}

	// DEBUG: CHECK THAT EVERY ROW IN Q IS EQUAL TO (r_j*s) ⊕ t_j
	print("CHECKING Q MATRIX ROWS: \n")
	for j := 0; j < m; j++ {
		for i := 0; i < k; i++ {
			q_idx := Q[j][i]
			r_j := receiver.SelectionBits[j]
			s_i := sender.S[i]
			t_j := T[j][i]

			if q_idx != (r_j*s_i)^t_j {
				print("r_j, s_i, t[j][i]: " + strconv.Itoa(int(r_j)) + " " + strconv.Itoa(int(s_i)) + " " + strconv.Itoa(int(t_j)) + "\n")
				print("q_idx: " + strconv.Itoa(int(q_idx)) + "\n")
				fmt.Println("ERROR: q_idx != (r_j * s_i) ^ t_j")
			}
		}
	}

	// The sender sends (y0_j, y1_j) for every 1 ≤ j ≤ m to the receiver, where y0_j = x0_j ⊕ H(j, q_j) and y1_j = x1_j ⊕ H(j, q_j ⊕ s).
	// The receiver then computes x^(r_j)_j = y^(rj)_j ⊕ H(j, t_j) for every 1 ≤ j ≤ m. Then outputs (x^(r_1)_1, ..., x^(r_m)_m).
	ByteCiphertexts := sender.MakeAndSendCiphertexts()
	result := receiver.DecryptCiphertexts(ByteCiphertexts)

	if len(result) != len(messages) {
		fmt.Println("Result length is not equal to messages length in OTExtensionProtocol")
	}

	print("Result: ")
	for _, b := range result {
		fmt.Printf("%d ", b)
	}
	fmt.Println()
}
