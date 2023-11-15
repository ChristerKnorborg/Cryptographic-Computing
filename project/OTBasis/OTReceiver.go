// OTReceiver.go
package project

// Import your ElGamal package
import (
	"crypto"
	"cryptographic-computing/project/elgamal"
	"math/big"
)

type OTReceiver struct {
	// Hold elGamal public parameters and secret key
	secretKeys []*big.Int
	choiceBits []int
}

func (receiver *OTReceiver) Init() {

}

func (receiver *OTReceiver) Choose(elGamal *elgamal.ElGamal, choices int) []*big.Int {

	// Generate secretkeys for each of the messages to be received
	for i := 0; i < choices; i++ {
		receiver.secretKeys[i] = elGamal.MakeSecretKey()
	}

	// Initialize a list of choices public keys to be sent to the OTSender
	publicKeys := make([]*big.Int, choices)

	for i := 0; i < choices; i++ {
		if i == receiver.choiceBits[i] {
			publicKeys[i] = elGamal.Gen(receiver.secretKeys[i]) // Generate the real public key corresponding to Alice's input x
		} else {
			publicKeys[i] = elGamal.OGen() // Generate 7 fake public keys using the oblivious version of Gen
		}
	}

	return publicKeys
}

func (receiver *OTReceiver) ReceiveData(elGamal *elgamal.ElGamal) {
	// Code to receive encrypted messages and public parameters
	xD := crypto.SHA256.New()
	if xD == nil {
		print("xD is nil")
	}
}

func (receiver *OTReceiver) DecryptMessage(choice int) string {
	// Decrypt the message based on the receiver's choice
	return ""
}
