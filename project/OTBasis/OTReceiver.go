// OTReceiver.go
package project

// Import your ElGamal package
import (
	"crypto"
	"cryptographic-computing/project/elgamal"
)

type OTReceiver struct {
	// State for receiver
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
