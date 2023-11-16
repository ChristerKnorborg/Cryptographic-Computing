// OTSender.go
package OTExtension

import (
	"crypto/rand"
	"cryptographic-computing/project/elgamal"
	"encoding/base64"
	"math/big"
)

type OTSender struct {
	m             int              // Number of messages to be sent
	secretKeys    []*big.Int       // Secret keys for each message to be received.
	PublicKeys    []*PublicKeyPair // Public keys received from the OTReceiver - one oblivious and one real for each message to be sent
	Messages      []*MessagePair   // Messages to be sent, each message consists of 2 messages M0 and M1.
	selectionBits []int            // Selection bits to invoke the κ×OTκ-functionality, where the OTSender plays the receiver
	s             string           // Random string s = (s_1, ... , s_k)
	seeds         []*big.Int       // Seed values received when invoking the κ×OTκ-functionality, where the OTSender plays the receiver and OTReceiver plays the sender.
	Q             [][]byte         // Bit matrix Q of size m × κ to be calculated in the OTExtension Phase
}

func (sender *OTSender) Init(Messages []*MessagePair) {

	sender.m = len(Messages)
	sender.Messages = Messages
}

// S choose a random string s = (s_1, ... , s_k)
func (sender *OTSender) ChooseRandomString(len int) {

	bytes := make([]byte, len) // Create a byte slice of the desired length

	// Read random bytes; note that rand.Read is from crypto/rand (cryptographically secure)
	_, err := rand.Read(bytes)
	if err != nil {
		panic("Error in ChooseRandomString: " + err.Error())
	}

	// Encode the bytes to a base64 string and return.
	// The length of the base64 string will be longer than len,
	// so we take a substring of the desired length
	sender.s = base64.RawURLEncoding.EncodeToString(bytes)[:len]
}

// Method for invoking the κ×OTκ-functionality, where the OTSender plays the receiver with random string s = (s_1, ... , s_k) as input,
// and OTReceiver plays the sender with inputs (k0_i, k1_i ) for every 1 ≤ i ≤ κ.
func (sender *OTSender) Choose(elGamal *elgamal.ElGamal) []*PublicKeyPair {

	k := len(sender.s)

	// Generate secretkeys for each of the messages to be received
	for i := 0; i < k; i++ {
		sender.secretKeys[i] = elGamal.MakeSecretKey()
	}

	// Initialize a list of choices public keys to be sent to the OTSender
	publicKeys := make([]*PublicKeyPair, k)

	for i := 0; i < k; i++ {
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
func (sender *OTSender) DecryptSeeds(ciphertextPairs []*CiphertextPair, elGamal *elgamal.ElGamal) []*big.Int {

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
	return plaintextSeeds
}

func (sender *OTSender) GenerateQMatrix(U [][]byte) {

	k := len(sender.s)
	m := len(U)

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
func (sender *OTSender) SendYValues() {

	k := len(sender.s)
	m := sender.m

	for j := 0; j < m; j++ {
		Q_row_j := sender.Q[j]
		y0_j := sender.Messages[j].Message0 ^ HashFunction(j, Q_row_j)
		y1_j := sender.Messages[j].Message1 ^ HashFunction(j, Q_row_j^sender.s)
	}

}
