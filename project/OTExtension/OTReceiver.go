// OTReceiver.go
package OTExtension

// Import your ElGamal package
import (
	"crypto/rand"
	"cryptographic-computing/project/elgamal"
	"math/big"
)

type OTReceiver struct {
	// Hold elGamal public parameters and secret key
	secretKeys    []*big.Int       // Secret keys for each message to be received.
	selectionBits []int            // Selection bits for each message to be received depending on if the receiver wants to learn M0 or M1 (Hidden for the OTSender).
	PublicKeys    []*PublicKeyPair // Public keys received from the OTSender when invoking the κ×OTκ-functionality, where the OTSender plays the receiver and OTReceiver plays the sender.
	seeds         []*seed          // Messages (seeds) to be sent, when invoking the κ×OTκ-functionality, where the OTSender plays the receiver and OTReceiver plays the sender.
	T             [][]byte         // Bit matrix T of size m × κ, after the κ×OTκ OT-functionality
}

// The receiver chooses k pairs of k-bit seeds {(k0_i , k1_i )} from i = 1 to k
func (receiver *OTReceiver) ChooseSeeds(k int) []*seed {

	seeds := make([]*seed, k)

	for i := 0; i < k; i++ {
		// Generate a k-bit random number for seed0
		seed0, err := rand.Int(rand.Reader, big.NewInt(1).Lsh(big.NewInt(1), uint(k)))
		if err != nil {
			panic("Error in ChooseSeeds for seed0: " + err.Error())
		}
		// Generate a k-bit random number for seed1
		seed1, err := rand.Int(rand.Reader, big.NewInt(1).Lsh(big.NewInt(1), uint(k)))
		if err != nil {
			panic("Error in ChooseSeeds for seed1: " + err.Error())
		}
		seeds[i] = &seed{
			seed0: seed0,
			seed1: seed1,
		}
	}
	return seeds
}

// Method for to receive Public keys, when the parties invoke the κ×OTκ-functionality, where the OTSender plays the receiver and OTReceiver plays the sender.
func (receiver *OTReceiver) ReceiveKeys(PublicKeys []*PublicKeyPair) {
	receiver.PublicKeys = PublicKeys
}

// Method to encrypt messages (seeds) when the parties invoke the κ×OTκ-functionality, where the OTSender plays the receiver and OTReceiver plays the sender.
func (receiver *OTReceiver) EncryptSeeds(elGamal *elgamal.ElGamal) []*CiphertextPair {

	k := len(receiver.seeds)

	ciphertexts := make([]*CiphertextPair, k)

	for i := 0; i < k; i++ {

		// Encrypt the messages using the public keys received from the OTReceiver
		msg0 := elGamal.Encrypt(receiver.PublicKeys[i].MessageKey0, receiver.seeds[i].seed0)
		msg1 := elGamal.Encrypt(receiver.PublicKeys[i].MessageKey1, receiver.seeds[i].seed1)

		// Store the encrypted messages in the ciphertext pair
		ciphertexts[i].Ciphertext0 = msg0
		ciphertexts[i].Ciphertext1 = msg1
	}

	return ciphertexts

}

// Method for generating the bit matrix T of size m × κ, after the κ×OTκ OT-functionality, where the OTSender plays the receiver and OTReceiver plays the sender.
// GenerateMatrixT generates the bit matrix T after the k×OTk functionality.
func (receiver *OTReceiver) GenerateMatrixT() {

	k := len(receiver.seeds)
	m := len(receiver.selectionBits)

	// Initialize the matrix T of size m × κ.
	T := make([][]byte, m) // m rows.
	for i := range T {
		T[i] = make([]byte, k) // k columns per row.
	}

	// Generate each column of T using the seeds.
	for i, seedPair := range receiver.seeds {
		// Generate a pseudo-random bit string of m bits.
		t_i, err := pseudoRandomGenerator(seedPair.seed0, m)
		if err != nil {
			panic("Error from pseudoRandomGenerator in GenerateMatrixT: " + err.Error())
		}

		// Assign each bit of t_i to the corresponding row.
		for j := 0; j < m; j++ {
			// Assuming t_i is a byte slice where each byte is either 0 or 1.
			T[j][i] = t_i[j]
		}
	}

	receiver.T = T
}

func (receiver *OTReceiver) GenerateAndSendUMatrix() [][]byte {

	k := len(receiver.seeds)
	m := len(receiver.selectionBits)

	// Initialize the matrix U of size m × κ.
	U := make([][]byte, m) // m rows.
	for i := range U {
		U[i] = make([]byte, k) // k columns per row.
	}

	// Generate each column of U: u^i = t^i ⊕ G(k1_i ) ⊕ r.
	for i := 0; i < k; i++ {

		// Generate a pseudo-random bit string of m bits.
		K1_PRG, err := pseudoRandomGenerator(receiver.seeds[i].seed1, m)

		if err != nil {
			panic("Error from pseudoRandomGenerator in GenerateAndSendUMatrix: " + err.Error())
		}

		col_i := make([]byte, m)
		for j := 0; j < m; j++ {

			col_i[j] = receiver.T[j][i] ^ K1_PRG[j] ^ byte(receiver.selectionBits[j])
		}
		U[i] = col_i
	}

	return U

}

func (receiver *OTReceiver) Init() {
	panic("implement me")
}

func (receiver *OTReceiver) DecryptMessage(ciphertextPairs []*CiphertextPair, elGamal *elgamal.ElGamal) []*big.Int {
	panic("implement me")
}

func (receiver *OTReceiver) ReceiveData(elGamal *elgamal.ElGamal) {
	panic("implement me")
}
