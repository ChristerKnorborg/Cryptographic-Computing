// OTSender.go
package OTBasic

import (
	"cryptographic-computing/project/elgamal"
	"cryptographic-computing/project/utils"
	"math/big"
)

type OTSender struct {
	PublicKeys []*utils.PublicKeyPair // Public keys received from the OTReceiver - one oblivious and one real for each message to be sent
	Messages   []*utils.MessagePair   // Messages to be sent, each message consists of 2 messages M0 and M1.
}

func (sender *OTSender) Init(Messages []*utils.MessagePair) {

	sender.Messages = make([]*utils.MessagePair, len(Messages))
	sender.PublicKeys = make([]*utils.PublicKeyPair, len(Messages))
	sender.Messages = Messages
}

func (sender *OTSender) ReceiveKeys(PublicKeys []*utils.PublicKeyPair) {
	sender.PublicKeys = PublicKeys
}

func (sender *OTSender) EncryptMessages(elGamal *elgamal.ElGamal) []*utils.CiphertextPair {

	ciphertexts := make([]*utils.CiphertextPair, len(sender.Messages))

	for i := 0; i < len(sender.Messages); i++ {

		ciphertexts[i] = &utils.CiphertextPair{} // Initialize a ciphertext pair for each message

		// Convert the messages to big integers
		msg0BigInt := new(big.Int)
		msg1BigInt := new(big.Int)
		msg0BigInt.SetBytes(sender.Messages[i].Message0)
		msg1BigInt.SetBytes(sender.Messages[i].Message1)

		// Encrypt the messages using the public keys received from the OTReceiver.
		// The sender is oblivious to which message is encrypted using which key.
		cipher0 := elGamal.Encrypt(msg0BigInt, sender.PublicKeys[i].MessageKey0)
		cipher1 := elGamal.Encrypt(msg1BigInt, sender.PublicKeys[i].MessageKey1)

		// Store the encrypted messages in the ciphertext pair
		ciphertexts[i].Ciphertext0 = cipher0
		ciphertexts[i].Ciphertext1 = cipher1
	}

	return ciphertexts

}
