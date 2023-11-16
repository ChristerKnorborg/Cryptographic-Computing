// OTSender.go
package OTBasic

import (
	"cryptographic-computing/project/elgamal"
)

type OTSender struct {
	PublicKeys []*PublicKeyPair // Public keys received from the OTReceiver - one oblivious and one real for each message to be sent
	Messages   []*MessagePair   // Messages to be sent, each message consists of 2 messages M0 and M1.
}

func (sender *OTSender) Init(Messages []*MessagePair, choices int) {
	for i := 0; i < choices; i++ {
		sender.Messages[i] = Messages[i]
	}
}

func (sender *OTSender) ReceiveKeys(PublicKeys []*PublicKeyPair) {
	sender.PublicKeys = PublicKeys
}

func (sender *OTSender) EncryptMessages(elGamal *elgamal.ElGamal) []*CiphertextPair {

	ciphertexts := make([]*CiphertextPair, len(sender.Messages))

	for i := 0; i < len(sender.Messages); i++ {

		// Encrypt the messages using the public keys received from the OTReceiver
		msg0 := elGamal.Encrypt(sender.PublicKeys[i].MessageKey0, sender.Messages[i].Message0)
		msg1 := elGamal.Encrypt(sender.PublicKeys[i].MessageKey1, sender.Messages[i].Message1)

		// Store the encrypted messages in the ciphertext pair
		ciphertexts[i].Ciphertext0 = msg0
		ciphertexts[i].Ciphertext1 = msg1
	}

	return ciphertexts

}

func (sender *OTSender) TransmitData(elGamal *elgamal.ElGamal) {
	// Code to transmit encrypted data and public parameters

}
