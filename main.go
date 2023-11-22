package main

import (
	"cryptographic-computing/project/OTExtension"
	"cryptographic-computing/project/elgamal"
	"fmt"
	"strconv"
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
	T := receiver.GenerateMatrixT() // REMEMBER TO REMOVE RETURN AFTER DEBUGGING
	U := receiver.GenerateAndSendMatrixU()
	Q := sender.GenerateMatrixQ(U) // REMEMBER TO REMOVE RETURN AFTER DEBUGGING

	m := len(selectionBits) // DEBUG
	// DEBUG: CHECK THAT EVERY COLUMN IN Q IS EQUAL TO (s_i * r) ⊕ t^i
	print("CHECKING Q MATRIX COLS: \n")
	for i := 0; i < k; i++ {
		s_i := sender.S[i : i+1]
		s_i_int, _ := strconv.Atoi(s_i)
		for j := 0; j < m; j++ {
			q_idx := Q[j][i]
			q_idx_int, _ := strconv.Atoi(q_idx)

			r_j := receiver.SelectionBits[j]
			t_i := T[j][i]
			t_i_int, _ := strconv.Atoi(t_i)

			if q_idx_int != (s_i_int*r_j)^t_i_int {
				G := q_idx_int ^ (s_i_int * r_j)
				print("s_i, u[j][i], G(): " + strconv.Itoa(s_i_int) + " " + U[j][i] + " " + strconv.Itoa(G) + "\n")
				print("s_i, r_j, t_i: " + strconv.Itoa(s_i_int) + " " + strconv.Itoa(r_j) + " " + strconv.Itoa(t_i_int) + "\n")
				print("q_idx: " + strconv.Itoa(q_idx_int) + "\n")
				fmt.Println("ERROR: q_idx != (s_i * r_j) ^ t_i")

			}
		}
	}

	// DEBUG: CHECK THAT EVERY ROW IN Q IS EQUAL TO (r_j*s) ⊕ t_j
	print("CHECKING Q MATRIX ROWS: \n")
	for j := 0; j < m; j++ {
		for i := 0; i < k; i++ {
			q_idx := Q[j][i]
			q_idx_int, _ := strconv.Atoi(q_idx)
			r_j := receiver.SelectionBits[j]
			r_j_int, _ := strconv.Atoi(strconv.Itoa(r_j))
			s_i := sender.S[i : i+1]
			s_i_int, _ := strconv.Atoi(s_i)
			t_j := T[j][i]
			t_j_int, _ := strconv.Atoi(t_j)

			if q_idx_int != (r_j_int*s_i_int)^t_j_int {
				print("r_j, s_i, t[j][i]: " + strconv.Itoa(r_j_int) + " " + strconv.Itoa(s_i_int) + " " + strconv.Itoa(t_j_int) + "\n")
				print("q_idx: " + strconv.Itoa(q_idx_int) + "\n")
				fmt.Println("ERROR: q_idx != (r_j * s_i) ^ t_j")

			}
		}
	}
	// Print seeds for sender
	print("Seeds sender: \n")
	for i := 0; i < k; i++ {
		print(sender.Seeds[i].String() + "\n")

	}

	// Print seeds for receiver
	print("Seeds receiver: \n")
	for i := 0; i < k; i++ {
		print(receiver.Seeds[i].Seed0.String() + "\n")
		print(receiver.Seeds[i].Seed1.String() + "\n")
	}

	// // The sender sends (y0_j, y1_j) for every 1 ≤ j ≤ m to the receiver, where y0_j = x0_j ⊕ H(j, q_j) and y1_j = x1_j ⊕ H(j, q_j ⊕ s).
	// // The receiver then computes x^(r_j)_j = y^(rj)_j ⊕ H(j, t_j) for every 1 ≤ j ≤ m. Then outputs (x^(r_1)_1, ..., x^(r_m)_m).
	// ByteCiphertexts := sender.MakeAndSendCiphertexts()
	// result := receiver.DecryptCiphertexts(ByteCiphertexts)

	// for _, b := range result {
	// 	fmt.Printf("%d ", b) // Decimal print of []byte result
	// }
}
