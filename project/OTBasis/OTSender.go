// OTSender.go
package project

import (
	"crypto"
	"cryptographic-computing/project/elgamal"
)

type OTSender struct {
	// Your ElGamal parameters and any other state
}

func (sender *OTSender) PrepareMessages(msg1, msg2 string) {

	crypto.SHA256.New()
	// Encrypt msg1 and msg2 using ElGamal
}

func (sender *OTSender) TransmitData(elGamal *elgamal.ElGamal) {
	// Code to transmit encrypted data and public parameters
}
