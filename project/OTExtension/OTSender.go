// OTSender.go
package OTExtension

import (
	"crypto/rand"
	"cryptographic-computing/project/elgamal"
	"cryptographic-computing/project/utils"
	"math/big"

	"github.com/hashicorp/vault/sdk/helper/xor"
)

type OTSender struct {
	messages   []*utils.MessagePair   // Sender S holds m pairs (x0_j, x1_j) of l-bit strings, for every 1 ≤ j ≤ m.
	l          int                    // Bit length of each message
	m          int                    // Number of messages to be sent
	k          int                    // Security parameter
	s          []byte                 // Random list of 0's and 1's: s = (s_1, ... , s_k).
	secretKeys []*big.Int             // Secret keys for each seed to be received and decrypted. (Used for k regular OTs)
	PublicKeys []*utils.PublicKeyPair // Public keys to be received from the OTReceiver - one oblivious and one real for each message to be sent
	seeds      []*big.Int             // Seed values to be received from the k regular OTs
	q          [][]byte               // Bit matrix Q of size m × κ to be calculated in the OTExtension Phase
}

func (sender *OTSender) Init(messages []*utils.MessagePair, k int, l int) {

	sender.l = l
	sender.m = len(messages)
	sender.k = k
	sender.messages = messages // Message pairs (x0_j, x1_j) of l-bit strings, for every 1 ≤ j ≤ m.

}

// S choose a random list of 0's and 1's of length k: s = (s_1, ... , s_k)
func (sender *OTSender) ChooseRandomS() {
	sender.s = make([]byte, sender.k) // Allocate space for the array

	for i := uint(0); i < uint(sender.k); i++ {
		randomBit, err := rand.Int(rand.Reader, big.NewInt(2)) // Generate a random bit
		if err != nil {
			panic("Error in ChooseRandomS: " + err.Error())
		}
		// Set 0 or 1 in the array
		sender.s[i] = byte(randomBit.Int64())
	}
}

// DEBUGGING METHOD. S choose a fixed list of 0's and 1's of length k: s = (s_1, ... , s_k)
func (sender *OTSender) ChooseFixedS() {

	fixedS := []byte{
		1, 1, 1, 1, // ... Define a fixed list of 0's and 1's up to sender.k elements
	}
	sender.s = make([]byte, sender.k) // Allocate space for the array

	for i := 0; i < sender.k; i++ {
		sender.s[i] = fixedS[i%len(fixedS)] // Make sure the number of fixed elements is at least sender.k
	}
}

// Method for the k regular OTs, where the OTSender plays the receiver with random string s = (s_1, ... , s_k) as input.
func (sender *OTSender) Choose(elGamal *elgamal.ElGamal) []*utils.PublicKeyPair {

	k := sender.k

	// Generate secretkeys for each of the messages to be received
	sender.secretKeys = make([]*big.Int, k)
	for i := 0; i < k; i++ {
		sender.secretKeys[i] = elGamal.MakeSecretKey()
	}

	// Make two public keys - one oblivious and one real for every message(seed) to be received.
	// The bits from the sender's string s chooses which public key to use for each message (seed).
	publicKeys := make([]*utils.PublicKeyPair, k)
	for i := 0; i < k; i++ {

		publicKeys[i] = &utils.PublicKeyPair{} // Initialize a new public key pair to store the keys for the current message

		if sender.s[i] == 0 {
			publicKeys[i].MessageKey0 = elGamal.Gen(sender.secretKeys[i])
			publicKeys[i].MessageKey1 = elGamal.OGen()
		} else if sender.s[i] == 1 {
			publicKeys[i].MessageKey0 = elGamal.OGen()
			publicKeys[i].MessageKey1 = elGamal.Gen(sender.secretKeys[i])
		} else {
			panic("Receiver string s bits are not 0 or 1 in Choose")
		}
	}
	return publicKeys
}

// Method to decrypt the Seeds (messages) sent by the OTReceiver, for the k regular OTs,
// where the OTSender plays the receiver, and OTReceiver plays the sender.
func (sender *OTSender) DecryptSeeds(ciphertextPairs []*utils.CiphertextPair, elGamal *elgamal.ElGamal) {

	// Initialize a list of Seeds to be decrypted
	plaintextSeeds := make([]*big.Int, len(ciphertextPairs))

	// Decrypt the message based on the receiver's bits from string s.
	for i := 0; i < sender.k; i++ {

		if sender.s[i] == 0 {
			plaintextSeeds[i] = elGamal.Decrypt(ciphertextPairs[i].Ciphertext0.C1, ciphertextPairs[i].Ciphertext0.C2, sender.secretKeys[i])
		} else if sender.s[i] == 1 {
			plaintextSeeds[i] = elGamal.Decrypt(ciphertextPairs[i].Ciphertext1.C1, ciphertextPairs[i].Ciphertext1.C2, sender.secretKeys[i])
		} else {
			panic("Receiver string s bits are not 0 or 1 in DecryptSeeds")
		}
	}
	sender.seeds = make([]*big.Int, sender.k)
	sender.seeds = plaintextSeeds
}

// Method for generating the bit matrix Q of size m × κ.
// Notice this method is inefficient, since it accesses the entire matrix T (e.g. m x k entries).
func (sender *OTSender) GenerateMatrixQ(U [][]byte) {

	k := sender.k
	m := sender.m

	// Initialize the matrix Q of size m × κ.
	Q := make([][]byte, m) // m rows.
	for i := range Q {
		Q[i] = make([]byte, k) // k cols per row.
	}

	// The OTSender defines q^i = (s_i · u^i) ⊕ G(k^(s_i)_i. Note that q^i = (s_i · r) ⊕ t^i)
	for i := 0; i < k; i++ {

		bitstring, err := utils.PseudoRandomGenerator(sender.seeds[i], m)
		for j := 0; j < m; j++ {
			if err != nil {
				panic("Error from pseudoRandomGenerator in GenerateQMatrix: " + err.Error())
			}

			// If the bit from string s is 0, q^i = ⊕ G(k^(0)_i
			if sender.s[i] == 0 {

				Q[j][i] = bitstring[j]

				// If the bit from string s is 1, q^i = u^i ⊕ G(k^(1)_i
			} else if sender.s[i] == 1 {

				Q[j][i] = U[j][i] ^ bitstring[j]

			} else {
				panic("Receiver s idx is not 0 or 1 in GenerateQMatrix")
			}
		}
	}
	sender.q = Q
}

// A more efficient method for generating the bit matrix Q of size m × κ.
// It generates the matrix Q row-wise for transposing afterwards.
func (sender *OTSender) GenerateMatrixQTranspose(U [][]byte) {

	k := sender.k
	m := sender.m

	// Initialize the matrix Q of size κ × m (transposed later).
	Q := make([][]byte, k) // m rows.
	for i := range Q {
		Q[i] = make([]byte, m) // k col per row.
	}

	// The OTSender defines q^i = (s_i · u^i) ⊕ G(k^(s_i)_i. Note that q^i = (s_i · r) ⊕ t^i)
	for i := 0; i < k; i++ {

		bitstring, err := utils.PseudoRandomGenerator(sender.seeds[i], m)
		if err != nil {
			panic("Error from pseudoRandomGenerator in GenerateQMatrix: " + err.Error())
		}

		// If the bit from string s is 0, q^i = ⊕ G(k^(0)_i
		if sender.s[i] == 0 {

			Q[i] = bitstring

			// If the bit from string s is 1, q^i = u^i ⊕ G(k^(1)_i
		} else if sender.s[i] == 1 {

			xorRes, err := xor.XORBytes(U[i], bitstring)
			if err != nil {
				panic("Error from XORBytes in GenerateMatrixQTranspose: " + err.Error())
			}
			Q[i] = xorRes

		} else {
			panic("Receiver S idx are not 0 or 1 in GenerateMatrixQTranspose")
		}
	}
	sender.q = utils.TransposeMatrix(Q) // Transpose the matrix Q
}

// An even more efficient method for generating the bit matrix Q of size m × κ.
// It generates the matrix Q row-wise for transposing afterwards using Eklundh's algorithm.
func (sender *OTSender) GenerateMatrixQEklundh(U [][]byte, multithreaded bool) {

	k := sender.k
	m := sender.m

	// Initialize the matrix Q of size m × κ (transposed later).
	Q := make([][]byte, k) // m rows.
	for i := range Q {
		Q[i] = make([]byte, m) // k columns per row.
	}

	// The OTSender defines q^i = (s_i · u^i) ⊕ G(k^(s_i)_i. Note that q^i = (s_i · r) ⊕ t^i)
	for i := 0; i < k; i++ {

		bitstring, err := utils.PseudoRandomGenerator(sender.seeds[i], m)
		if err != nil {
			panic("Error from pseudoRandomGenerator in GenerateQMatrixEklundh: " + err.Error())
		}

		// If the bit from string s is 0, q^i = ⊕ G(k^(0)_i
		if sender.s[i] == 0 {

			Q[i] = bitstring

			// If the bit from string s is 1, q^i = u^i ⊕ G(k^(1)_i
		} else if sender.s[i] == 1 {

			xorRes, err := xor.XORBytes(U[i], bitstring)
			if err != nil {
				panic("Error from XORBytes in GenerateQMatrixEklundh: " + err.Error())
			}
			Q[i] = xorRes

		} else {
			panic("Receiver S idx are not 0 or 1 in GenerateQMatrix")
		}
	}
	sender.q = utils.EklundhTranspose(Q, multithreaded) // Transpose the matrix Q using Eklundh's algorithm
}

// Method for generating the ciphertexts to be sent to the OTReceiver.
// The OTSender sends m ciphertext pairs (y0_j, y1_j) of l-bit strings, for every 1 ≤ j ≤ m,
// where y0_j = x0_j ⊕ H(q_j) and y1_j = x1_j ⊕ H(q_j ⊕ s).
func (sender *OTSender) MakeAndSendCiphertexts() []*utils.ByteCiphertextPair {

	m := sender.m
	l := sender.l

	ByteCiphertextPairs := make([]*utils.ByteCiphertextPair, m)

	for j := 0; j < m; j++ {
		x0_j := sender.messages[j].Message0
		x1_j := sender.messages[j].Message1

		xor_res, err := xor.XORBytes(sender.q[j], sender.s) // XOR the j'th row of Q with the Sender's string s.
		if err != nil {
			panic("Error from XORBytes in MakeAndSendCiphertexts: " + err.Error())
		}

		hash0 := utils.Hash(sender.q[j], l) // Generate hash of length l from the j'th row of Q.
		hash1 := utils.Hash(xor_res, l)     // Generate hash of length l from the XOR of the j'th row of Q and the Sender's string s.

		y0_j, err1 := xor.XORBytes(x0_j, hash0)
		y1_j, err2 := xor.XORBytes(x1_j, hash1)

		if err1 != nil || err2 != nil {
			panic("Error from XORBytes in MakeAndSendCiphertexts: " + err1.Error() + " " + err2.Error())
		}

		ByteCiphertextPairs[j] = &utils.ByteCiphertextPair{Y0: y0_j, Y1: y1_j}
	}
	return ByteCiphertextPairs
}
