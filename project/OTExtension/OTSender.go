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
	messages   []*utils.MessagePair   // Sender S holds m pairs (x0_j, x1_j of l-bit strings, for every 1 ≤ j ≤ m.
	l          int                    // Bit length of each message
	m          int                    // Number of messages to be sent
	k          int                    // Security parameter
	s          []uint8                // Random list of 0's and 1's: s = (s_1, ... , s_k). Notice, we use uint8 instead of string due to XOR operations in GO.
	secretKeys []*big.Int             // Secret keys for each seeds to be received.
	PublicKeys []*utils.PublicKeyPair // Public keys received from the OTReceiver - one oblivious and one real for each message to be sent
	seeds      []*big.Int             // Seed values received when invoking the κ×OTκ-functionality, where the OTSender plays the receiver and OTReceiver plays the sender.
	q          [][]uint8              // Bit matrix Q of size m × κ to be calculated in the OTExtension Phase
}

func (sender *OTSender) Init(messages []*utils.MessagePair, securityParameter int, l int) {

	sender.l = l
	sender.m = len(messages)
	sender.k = securityParameter
	sender.messages = messages // Message pairs (x0_j, x1_j) of l-bit strings, for every 1 ≤ j ≤ m.

}

// S choose a random list of 0's and 1's of length k: s = (s_1, ... , s_k)
func (sender *OTSender) ChooseRandomK() {
	sender.s = make([]uint8, sender.k) // Allocate space for the array

	for i := uint(0); i < uint(sender.k); i++ {
		randomBit, err := rand.Int(rand.Reader, big.NewInt(2))
		if err != nil {
			panic("Error in ChooseRandomUint8Array: " + err.Error())
		}
		// Set 0 or 1 in the array
		sender.s[i] = uint8(randomBit.Int64())
	}
}

// Method for invoking the κ×OTκ-functionality, where the OTSender plays the receiver with random string s = (s_1, ... , s_k) as input,
// and OTReceiver plays the sender with inputs (k0_i, k1_i ) for every 1 ≤ i ≤ κ.
func (sender *OTSender) Choose(elGamal *elgamal.ElGamal) []*utils.PublicKeyPair {

	k := sender.k

	sender.secretKeys = make([]*big.Int, k)
	// Generate secretkeys for each of the messages to be received
	for i := 0; i < k; i++ {
		sender.secretKeys[i] = elGamal.MakeSecretKey()
	}

	// Initialize a list of choices public keys to be sent to the OTSender
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
			panic("Receiver choice bits are not 0 or 1 in Choose")
		}
	}

	return publicKeys
}

// Method to decrypt the Seeds (messages) sent by the OTReceiver, when invoking the κ×OTκ-functionality, where the OTSender plays the receiver, and OTReceiver plays the sender.
func (sender *OTSender) DecryptSeeds(ciphertextPairs []*utils.CiphertextPair, elGamal *elgamal.ElGamal) {

	// Initialize a list of Seeds to be decrypted
	plaintextSeeds := make([]*big.Int, len(ciphertextPairs))

	// Decrypt the message based on the receiver's choice bits
	for i := 0; i < sender.k; i++ {

		if sender.s[i] == 0 {
			plaintextSeeds[i] = elGamal.Decrypt(ciphertextPairs[i].Ciphertext0.C1, ciphertextPairs[i].Ciphertext0.C2, sender.secretKeys[i])
		} else if sender.s[i] == 1 {
			plaintextSeeds[i] = elGamal.Decrypt(ciphertextPairs[i].Ciphertext1.C1, ciphertextPairs[i].Ciphertext1.C2, sender.secretKeys[i])
		} else {
			panic("Receiver choice bits are not 0 or 1 in DecryptMessage")
		}
	}
	sender.seeds = make([]*big.Int, sender.k)
	sender.seeds = plaintextSeeds
}

func (sender *OTSender) GenerateMatrixQ(U [][]uint8) {

	k := sender.k
	m := sender.m

	// Initialize the matrix Q of size m × κ.
	Q := make([][]uint8, m) // m rows.
	for i := range Q {
		Q[i] = make([]uint8, k) // k columns per row.
	}

	// The OTSender defines q^i = (s_i · u^i) ⊕ G(k^(s_i)_i. Note that q^i = (s_i · r) ⊕ t^i)
	for i := 0; i < k; i++ {

		bitstring, err := utils.PseudoRandomGenerator(sender.seeds[i], m)
		for j := 0; j < m; j++ {
			if err != nil {
				panic("Error from pseudoRandomGenerator in GenerateQMatrix: " + err.Error())
			}

			// The reduction where G(k0_i) = G(k1_i) if the bit from string s is 0, can be seen on page 15 in the ALSZ paper.
			if sender.s[i] == 0 {

				Q[j][i] = bitstring[j]

			} else if sender.s[i] == 1 {

				Q[j][i] = U[j][i] ^ bitstring[j]

			} else {
				panic("Receiver S idx are not 0 or 1 in GenerateQMatrix")
			}
		}
	}
	sender.q = Q
}

func (sender *OTSender) GenerateMatrixQEklundh(U [][]uint8, multithreaded bool) {

	k := sender.k
	m := sender.m

	// Initialize the matrix Q of size m × κ.
	Q := make([][]uint8, k) // m rows.
	for i := range Q {
		Q[i] = make([]uint8, m) // k columns per row.
	}

	// The OTSender defines q^i = (s_i · u^i) ⊕ G(k^(s_i)_i. Note that q^i = (s_i · r) ⊕ t^i)
	for i := 0; i < k; i++ {

		bitstring, err := utils.PseudoRandomGenerator(sender.seeds[i], m)
		if err != nil {
			panic("Error from pseudoRandomGenerator in GenerateQMatrix: " + err.Error())
		}

		// The reduction where G(k0_i) = G(k1_i) if the bit from string s is 0, can be seen on page 15 in the ALSZ paper.
		if sender.s[i] == 0 {

			Q[i] = bitstring

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
	sender.q = utils.EklundhTranspose(Q, multithreaded)
}

func (sender *OTSender) GenerateMatrixQTranspose(U [][]uint8) {

	k := sender.k
	m := sender.m

	// Initialize the matrix Q of size m × κ.
	Q := make([][]uint8, k) // m rows.
	for i := range Q {
		Q[i] = make([]uint8, m) // k columns per row.
	}

	// The OTSender defines q^i = (s_i · u^i) ⊕ G(k^(s_i)_i. Note that q^i = (s_i · r) ⊕ t^i)
	for i := 0; i < k; i++ {

		bitstring, err := utils.PseudoRandomGenerator(sender.seeds[i], m)
		if err != nil {
			panic("Error from pseudoRandomGenerator in GenerateQMatrix: " + err.Error())
		}

		// The reduction where G(k0_i) = G(k1_i) if the bit from string s is 0, can be seen on page 15 in the ALSZ paper.
		if sender.s[i] == 0 {

			Q[i] = bitstring

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
	sender.q = utils.TransposeMatrix(Q)
}

func (sender *OTSender) MakeAndSendCiphertexts() []*utils.ByteCiphertextPair {

	m := sender.m
	l := sender.l

	// Initialize a list of ciphertext pairs to be sent to the OTReceiver
	ByteCiphertextPairs := make([]*utils.ByteCiphertextPair, m)

	for j := 0; j < m; j++ {
		x0_j := sender.messages[j].Message0
		x1_j := sender.messages[j].Message1

		xor_res, err := xor.XORBytes(sender.q[j], sender.s)
		if err != nil {
			panic("Error from XORBytes in MakeAndSendCiphertexts: " + err.Error())
		}

		hash0 := utils.Hash(sender.q[j], l)
		hash1 := utils.Hash(xor_res, l)

		y0_j, err1 := xor.XORBytes(x0_j, hash0)
		y1_j, err2 := xor.XORBytes(x1_j, hash1)

		if err1 != nil || err2 != nil {
			panic("Error from XORBytes in MakeAndSendCiphertexts: " + err1.Error() + " " + err2.Error())
		}

		ByteCiphertextPairs[j] = &utils.ByteCiphertextPair{Y0: y0_j, Y1: y1_j}
	}

	return ByteCiphertextPairs
}
