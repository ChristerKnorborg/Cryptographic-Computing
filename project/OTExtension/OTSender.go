// OTSender.go
package OTExtension

import (
	"crypto/rand"
	"cryptographic-computing/project/elgamal"
	"encoding/base64"
	"math/big"
)

type OTSender struct {
	messages      []*MessagePair   // Sender S holds m pairs (x0_j, x1_j of l-bit strings, for every 1 ≤ j ≤ m.
	m             int              // Number of messages to be sent
	k             int              // Security parameter
	s             string           // Random string s = (s_1, ... , s_k)
	secretKeys    []*big.Int       // Secret keys for each message to be received.
	PublicKeys    []*PublicKeyPair // Public keys received from the OTReceiver - one oblivious and one real for each message to be sent
	selectionBits []int            // Selection bits to invoke the κ×OTκ-functionality, where the OTSender plays the receiver
	seeds         []*big.Int       // Seed values received when invoking the κ×OTκ-functionality, where the OTSender plays the receiver and OTReceiver plays the sender.
	Q             [][]byte         // Bit matrix Q of size m × κ to be calculated in the OTExtension Phase
}

func (sender *OTSender) Init(messages []*MessagePair, securityParameter int, selectionBits []int) {

	sender.m = len(messages)
	sender.k = securityParameter
	sender.messages = messages // Message pairs (x0_j, x1_j) of l-bit strings, for every 1 ≤ j ≤ m.
	sender.selectionBits = selectionBits

	if len(selectionBits) != sender.k {
		panic("Length of selection bits and messages are not equal")
	}

}

// S choose a random string s = (s_1, ... , s_k)
func (sender *OTSender) ChooseRandomString() {

	bytes := make([]byte, sender.k) // Create a byte slice of the desired length

	// Read random bytes; note that rand.Read is from crypto/rand (cryptographically secure)
	_, err := rand.Read(bytes)
	if err != nil {
		panic("Error in ChooseRandomString: " + err.Error())
	}

	// Encode the bytes to a base64 string and return.
	// The length of the base64 string will be longer than len,
	// so we take a substring of the desired length
	sender.s = base64.RawURLEncoding.EncodeToString(bytes)[:sender.k]
}

// Method for invoking the κ×OTκ-functionality, where the OTSender plays the receiver with random string s = (s_1, ... , s_k) as input,
// and OTReceiver plays the sender with inputs (k0_i, k1_i ) for every 1 ≤ i ≤ κ.
func (sender *OTSender) Choose(elGamal *elgamal.ElGamal) []*PublicKeyPair {

	k := sender.k

	sender.secretKeys = make([]*big.Int, k)
	// Generate secretkeys for each of the messages to be received
	for i := 0; i < k; i++ {
		sender.secretKeys[i] = elGamal.MakeSecretKey()
	}

	// Initialize a list of choices public keys to be sent to the OTSender
	publicKeys := make([]*PublicKeyPair, k)

	for i := 0; i < k; i++ {

		publicKeys[i] = &PublicKeyPair{} // Initialize a new public key pair to store the keys for the current message

		if sender.selectionBits[i] == 0 {
			publicKeys[i].MessageKey0 = elGamal.Gen(sender.secretKeys[i])
			publicKeys[i].MessageKey1 = elGamal.OGen()
		} else if sender.selectionBits[i] == 1 {
			publicKeys[i].MessageKey0 = elGamal.OGen()
			publicKeys[i].MessageKey1 = elGamal.Gen(sender.secretKeys[i])
		} else {
			panic("Receiver choice bits are not 0 or 1 in Choose")
		}
	}

	return publicKeys
}

// Method to decrypt the seeds (messages) sent by the OTReceiver, when invoking the κ×OTκ-functionality, where the OTSender plays the receiver, and OTReceiver plays the sender.
func (sender *OTSender) DecryptSeeds(ciphertextPairs []*CiphertextPair, elGamal *elgamal.ElGamal) {

	// Initialize a list of seeds to be decrypted
	plaintextSeeds := make([]*big.Int, len(ciphertextPairs))

	// Decrypt the message based on the receiver's choice bits
	for i := 0; i < len(ciphertextPairs); i++ {
		if sender.selectionBits[i] == 0 {
			plaintextSeeds[i] = elGamal.Decrypt(ciphertextPairs[i].Ciphertext0.C1, ciphertextPairs[i].Ciphertext0.C2, sender.secretKeys[i])
		} else if sender.selectionBits[i] == 1 {
			plaintextSeeds[i] = elGamal.Decrypt(ciphertextPairs[i].Ciphertext1.C1, ciphertextPairs[i].Ciphertext1.C2, sender.secretKeys[i])
		} else {
			panic("Receiver choice bits are not 0 or 1 in DecryptMessage")
		}
	}
	sender.seeds = make([]*big.Int, len(plaintextSeeds))
	sender.seeds = plaintextSeeds

	print("Sender seeds len: ", len(sender.seeds), "\n")
	print("Sender seeds: ", sender.seeds, "\n")
	for i := 0; i < len(sender.seeds); i++ {
		print("Sender seeds ", i, " as string: ", sender.seeds[i].String(), "\n")
	}
}

func (sender *OTSender) GenerateQMatrix(U [][]byte) {

	print("U: ")
	PrintMatrix(U)

	k := sender.k
	m := sender.m

	// Initialize the matrix Q of size m × κ.
	Q := make([][]byte, m) // m rows.
	for i := range Q {
		Q[i] = make([]byte, k) // k columns per row.
	}

	// The OTSender defines q^i = (s_i · u^i) ⊕ G(k^(s_i)_i. Note that q^i = (s_i · r) ⊕ t^i)
	for i := 0; i < k; i++ {
		G, err := pseudoRandomGenerator(sender.seeds[sender.selectionBits[i]], m)
		for j := 0; j < m; j++ {
			if err != nil {
				panic("Error from pseudoRandomGenerator in GenerateQMatrix: " + err.Error())
			}
			Q[j][i] = sender.s[i] ^ U[j][i] ^ G[i]
		}
	}
	sender.Q = Q
}

// The OTSender sends (y0_j , y1_j ) for every 1 ≤ j ≤ m, where y0_j = x0_j ⊕ H(j, q_j) and y1_j = x1_j ⊕ H(j, q_j ⊕ s)
// func (sender *OTSender) SendYValues() {

// 	panic("Not implemented")

// 	k := sender.k
// 	m := sender.m

// 	for j := 0; j < m; j++ {
// 		Q_row_j := sender.Q[j]
// 		y0_j := sender.Messages[j].Message0 ^ HashFunction(j, Q_row_j)
// 		y1_j := sender.Messages[j].Message1 ^ HashFunction(j, Q_row_j^sender.s)
// 	}

// }
