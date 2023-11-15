// OTReceiver.go
package project

// Import your ElGamal package
import (
	"cryptographic-computing/project/elgamal"
	"math/big"
)

type OTReceiver struct {
	// Hold elGamal public parameters and secret key
	secretKeys []*big.Int // Secret keys for each message to be received.
	choiceBits []int      // Choice bits for each message to be received depending on if the receiver wants to learn M0 or M1 (Hidden for the OTSender)
}

func (receiver *OTReceiver) Init() {
	panic("implement me")
}

func (receiver *OTReceiver) Choose(elGamal *elgamal.ElGamal, choices int) []*PublicKeyPair {

	// Generate secretkeys for each of the messages to be received
	for i := 0; i < choices; i++ {
		receiver.secretKeys[i] = elGamal.MakeSecretKey()
	}

	// Initialize a list of choices public keys to be sent to the OTSender
	publicKeys := make([]*PublicKeyPair, choices)

	for i := 0; i < choices; i++ {
		if receiver.choiceBits[i] == 0 {
			publicKeys[i].MessageKey0 = elGamal.Gen(receiver.secretKeys[i])
			publicKeys[i].MessageKey1 = elGamal.OGen() // Generate the real public key corresponding to Alice's input x
		} else if receiver.choiceBits[i] == 1 {
			publicKeys[i].MessageKey0 = elGamal.OGen() // Generate 7 fake public keys using the oblivious version of Gen
			publicKeys[i].MessageKey1 = elGamal.Gen(receiver.secretKeys[i])
		} else {
			panic("Receiver choice bits are not 0 or 1 in Choose")
		}
	}

	return publicKeys
}

func (receiver *OTReceiver) DecryptMessage(ciphertextPairs []*CiphertextPair, elGamal *elgamal.ElGamal) []*big.Int {

	// Initialize a list of plaintexts to be decrypted
	plaintexts := make([]*big.Int, len(ciphertextPairs))

	// Decrypt the message based on the receiver's choice bits
	for i := 0; i < len(ciphertextPairs); i++ {
		if receiver.choiceBits[i] == 0 {
			plaintexts[i] = elGamal.Decrypt(ciphertextPairs[i].Ciphertext0.C1, ciphertextPairs[i].Ciphertext0.C2, receiver.secretKeys[i])
		} else if receiver.choiceBits[i] == 1 {
			plaintexts[i] = elGamal.Decrypt(ciphertextPairs[i].Ciphertext1.C1, ciphertextPairs[i].Ciphertext1.C2, receiver.secretKeys[i])
		} else {
			panic("Receiver choice bits are not 0 or 1 in DecryptMessage")
		}
	}
	return plaintexts

}

func (receiver *OTReceiver) ReceiveData(elGamal *elgamal.ElGamal) {
	panic("implement me")
}
