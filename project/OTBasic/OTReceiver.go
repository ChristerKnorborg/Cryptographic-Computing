// OTReceiver.go
package OTBasic

// Import your ElGamal package
import (
	"cryptographic-computing/project/elgamal"
	"math/big"
)

type OTReceiver struct {
	// Hold elGamal public parameters and secret key
	secretKeys    []*big.Int // Secret keys for each message to be received.
	selectionBits []int      // Choice bits for each message to be received depending on if the receiver wants to learn M0 or M1 (Hidden for the OTSender)
}

func (receiver *OTReceiver) Init(selectionBits []int) {
	receiver.selectionBits = selectionBits
}

func (receiver *OTReceiver) Choose(num_selections int, elGamal *elgamal.ElGamal) []*PublicKeyPair {

	receiver.secretKeys = make([]*big.Int, num_selections)

	// Generate secretkeys for each of the messages to be received
	for i := 0; i < num_selections; i++ {
		receiver.secretKeys[i] = elGamal.MakeSecretKey()
	}

	// Initialize a list of num_selections public keys to be sent to the OTSender
	publicKeys := make([]*PublicKeyPair, num_selections)

	for i := 0; i < num_selections; i++ {

		publicKeys[i] = &PublicKeyPair{} // Initialize a public key pair for each message

		if receiver.selectionBits[i] == 0 {
			publicKeys[i].MessageKey0 = elGamal.Gen(receiver.secretKeys[i])
			publicKeys[i].MessageKey1 = elGamal.OGen()
		} else if receiver.selectionBits[i] == 1 {
			publicKeys[i].MessageKey0 = elGamal.OGen()
			publicKeys[i].MessageKey1 = elGamal.Gen(receiver.secretKeys[i])
		} else {
			panic("Receiver choice bits are not 0 or 1 in Choose")
		}
	}

	return publicKeys
}

func (receiver *OTReceiver) DecryptMessage(ciphertextPairs []*CiphertextPair, l int, elGamal *elgamal.ElGamal) [][]byte {

	// Initialize a list of plaintexts to be decrypted
	plaintexts := make([][]byte, len(ciphertextPairs))

	// Decrypt the message based on the receiver's choice bits
	for i, pair := range ciphertextPairs {
		var plaintext *big.Int
		if receiver.selectionBits[i] == 0 {
			plaintext = elGamal.Decrypt(pair.Ciphertext0.C1, pair.Ciphertext0.C2, receiver.secretKeys[i])
		} else if receiver.selectionBits[i] == 1 {
			plaintext = elGamal.Decrypt(pair.Ciphertext1.C1, pair.Ciphertext1.C2, receiver.secretKeys[i])
		} else {
			panic("Receiver choice bits are not 0 or 1 in DecryptMessage")
		}

		bytePlaintext := plaintext.Bytes() // Convert *big.Int to []byte

		// Ensure the byte slice is exactly 'l' bytes long
		if len(bytePlaintext) < l {
			// Padding at the end
			padding := make([]byte, l-len(bytePlaintext))
			bytePlaintext = append(bytePlaintext, padding...)
		} else if len(bytePlaintext) > l {
			// Truncate if longer
			bytePlaintext = bytePlaintext[:l]
		}

		plaintexts[i] = bytePlaintext
	}
	return plaintexts
}
