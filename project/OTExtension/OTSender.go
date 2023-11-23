// OTSender.go
package OTExtension

import (
	"crypto/rand"
	"cryptographic-computing/project/elgamal"
	"math/big"
	"strconv"
	"strings"

	"github.com/hashicorp/vault/sdk/helper/xor"
)

type OTSender struct {
	messages   []*MessagePair   // Sender S holds m pairs (x0_j, x1_j of l-bit strings, for every 1 ≤ j ≤ m.
	l          int              // Bit length of each message
	m          int              // Number of messages to be sent
	k          int              // Security parameter
	s          string           // Random string s = (s_1, ... , s_k).
	secretKeys []*big.Int       // Secret keys for each message to be received.
	PublicKeys []*PublicKeyPair // Public keys received from the OTReceiver - one oblivious and one real for each message to be sent
	seeds      []*big.Int       // Seed values received when invoking the κ×OTκ-functionality, where the OTSender plays the receiver and OTReceiver plays the sender.
	Q          [][]string       // Bit matrix Q of size m × κ to be calculated in the OTExtension Phase
}

func (sender *OTSender) Init(messages []*MessagePair, securityParameter int, l int) {

	sender.l = l
	sender.m = len(messages)
	sender.k = securityParameter
	sender.messages = messages // Message pairs (x0_j, x1_j) of l-bit strings, for every 1 ≤ j ≤ m.

}

// S choose a random string s = (s_1, ... , s_k)
func (sender *OTSender) ChooseRandomString() {
	var stringBuilder strings.Builder
	stringBuilder.Grow(sender.k) // Pre-allocate space for efficiency

	// Generate each bit individually and append to the string builder
	for i := 0; i < sender.k; i++ {
		randomBit, err := rand.Int(rand.Reader, big.NewInt(2))
		if err != nil {
			panic("Error in ChooseRandomString: " + err.Error())
		}
		// Append '0' or '1' to the string
		if randomBit.Int64() == 0 {
			stringBuilder.WriteByte('0')
		} else {
			stringBuilder.WriteByte('1')
		}
	}

	// Set the generated string
	sender.s = stringBuilder.String()
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
		s_idx, err := strconv.Atoi(sender.s[i : i+1])

		if err != nil {
			panic("Error from Atoi in Choose: " + err.Error())
		}

		if s_idx == 0 {
			publicKeys[i].MessageKey0 = elGamal.Gen(sender.secretKeys[i])
			publicKeys[i].MessageKey1 = elGamal.OGen()
		} else if s_idx == 1 {
			publicKeys[i].MessageKey0 = elGamal.OGen()
			publicKeys[i].MessageKey1 = elGamal.Gen(sender.secretKeys[i])
		} else {
			panic("Receiver choice bits are not 0 or 1 in Choose")
		}
	}

	return publicKeys
}

// Method to decrypt the Seeds (messages) sent by the OTReceiver, when invoking the κ×OTκ-functionality, where the OTSender plays the receiver, and OTReceiver plays the sender.
func (sender *OTSender) DecryptSeeds(ciphertextPairs []*CiphertextPair, elGamal *elgamal.ElGamal) {

	// Initialize a list of Seeds to be decrypted
	plaintextSeeds := make([]*big.Int, len(ciphertextPairs))

	// Decrypt the message based on the receiver's choice bits
	for i := 0; i < sender.k; i++ {
		s_idx, err := strconv.Atoi(sender.s[i : i+1])
		if err != nil {
			panic("Error from Atoi in DecryptMessage: " + err.Error())
		}

		if s_idx == 0 {
			plaintextSeeds[i] = elGamal.Decrypt(ciphertextPairs[i].Ciphertext0.C1, ciphertextPairs[i].Ciphertext0.C2, sender.secretKeys[i])
		} else if s_idx == 1 {
			plaintextSeeds[i] = elGamal.Decrypt(ciphertextPairs[i].Ciphertext1.C1, ciphertextPairs[i].Ciphertext1.C2, sender.secretKeys[i])
		} else {
			panic("Receiver choice bits are not 0 or 1 in DecryptMessage")
		}
	}
	sender.seeds = make([]*big.Int, sender.k)
	sender.seeds = plaintextSeeds
}

func (sender *OTSender) GenerateMatrixQ(U [][]string) {

	k := sender.k
	m := sender.m

	// Initialize the matrix Q of size m × κ.
	Q := make([][]string, m) // m rows.
	for i := range Q {
		Q[i] = make([]string, k) // k columns per row.
	}

	// The OTSender defines q^i = (s_i · u^i) ⊕ G(k^(s_i)_i. Note that q^i = (s_i · r) ⊕ t^i)
	for i := 0; i < k; i++ {
		s_idx, err := strconv.Atoi(sender.s[i : i+1])
		if err != nil {
			panic("Error from Atoi in GenerateQMatrix on string s: " + err.Error())
		}
		bitstring, err := pseudoRandomGenerator(sender.seeds[i], m)
		for j := 0; j < m; j++ {
			if err != nil {
				panic("Error from pseudoRandomGenerator in GenerateQMatrix: " + err.Error())
			}

			// The reduction where G(k0_i) = G(k1_i) if the bit from string s is 0, can be seen on page 15 in the ALSZ paper.
			if s_idx == 0 { // If s[i] = 1
				Q[j][i] = bitstring[j : j+1]

			} else if s_idx == 1 { // If s[i] = 1
				bitstring_idx := bitstring[j : j+1]
				UEntry, err := strconv.Atoi(U[j][i])

				if err != nil {
					panic("Error from Atoi in GenerateQMatrix: " + err.Error())
				}

				xor, err := XOR(UEntry, bitstring_idx)
				if err != nil {
					panic("Error from XOR in GenerateQMatrix: " + err.Error())
				}
				Q[j][i] = xor

			} else {
				panic("Receiver S idx are not 0 or 1 in GenerateQMatrix")
			}
		}
	}
	sender.Q = Q
}

func (sender *OTSender) MakeAndSendCiphertexts() []*ByteCiphertextPair {

	k := sender.k
	m := sender.m
	l := sender.l

	// Initialize a list of ciphertext pairs to be sent to the OTReceiver
	ByteCiphertextPairs := make([]*ByteCiphertextPair, m)

	for j := 0; j < m; j++ {
		x0_j := sender.messages[j].Message0
		x1_j := sender.messages[j].Message1

		string_xor := ""
		q_row := ""
		for i := 0; i < k; i++ {
			string_idx := sender.s[i : i+1]
			q_idx := sender.Q[j][i]

			xor_char, err := XOR(string_idx, q_idx)

			if err != nil {
				panic("Error from XOR in MakeAndSendCiphertexts: " + err.Error())
			}
			string_xor += xor_char
			q_row += q_idx
		}

		hash0 := Hash([]byte(q_row), l)
		hash1 := Hash([]byte(string_xor), l)

		y0_j, err1 := xor.XORBytes(x0_j, hash0)
		y1_j, err2 := xor.XORBytes(x1_j, hash1)

		if err1 != nil || err2 != nil {
			panic("Error from XORBytes in MakeAndSendCiphertexts: " + err1.Error() + " " + err2.Error())
		}

		ByteCiphertextPairs[j] = &ByteCiphertextPair{y0: y0_j, y1: y1_j}
	}

	return ByteCiphertextPairs
}
