// OTSender.go
package OTExtension

import (
	"crypto/rand"
	"cryptographic-computing/project/elgamal"
	"encoding/base64"
	"math/big"
)

type OTSender struct {
	secretKeys    []*big.Int       // Secret keys for each message to be received.
	PublicKeys    []*PublicKeyPair // Public keys received from the OTReceiver - one oblivious and one real for each message to be sent
	Messages      []*MessagePair   // Messages to be sent, each message consists of 2 messages M0 and M1.
	selectionBits []int            // Selection bits to invoke the κ×OTκ-functionality, where the OTSender plays the receiver
	s             string           // Random string s = (s_1, ... , s_k)
	seeds         []*big.Int       // Seed values received when invoking the κ×OTκ-functionality, where the OTSender plays the receiver and OTReceiver plays the sender.
}

func (sender *OTSender) Init(Messages []*MessagePair, choices int) {
	for i := 0; i < choices; i++ {
		sender.Messages[i] = Messages[i]
	}
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

func (sender *OTSender) TransmitData(elGamal *elgamal.ElGamal) {
	// Code to transmit encrypted data and public parameters

}
