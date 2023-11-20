// OTReceiver.go
package OTExtension

// Import your ElGamal package
import (
	"crypto/rand"
	"cryptographic-computing/project/elgamal"
	"math/big"
	"strconv"
)

type OTReceiver struct {
	m             int              // Number of messages to be received
	k             int              // Security parameter
	selectionBits []int            // Receiver R holds m selection bits r = (r_1, ..., r_m).
	seeds         []*seed          // Messages (seeds) to be sent, when invoking the κ×OTκ-functionality, where the OTSender plays the receiver and OTReceiver plays the sender.
	secretKeys    []*big.Int       // Secret keys for each message to be received.
	PublicKeys    []*PublicKeyPair // Public keys received from the OTSender when invoking the κ×OTκ-functionality, where the OTSender plays the receiver and OTReceiver plays the sender.
	T             [][]string       // Bit matrix T of size m × κ, after the κ×OTκ OT-functionality
}

func (receiver *OTReceiver) Init(selectionBits []int, securityParameter int) {

	receiver.m = len(selectionBits)
	receiver.selectionBits = selectionBits
	receiver.k = securityParameter

}

// The receiver chooses k pairs of k-bit seeds {(k0_i , k1_i )} from i = 1 to k
func (receiver *OTReceiver) ChooseSeeds() {

	k := receiver.k

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
	receiver.seeds = seeds
}

// Method for to receive Public keys, when the parties invoke the κ×OTκ-functionality, where the OTSender plays the receiver and OTReceiver plays the sender.
func (receiver *OTReceiver) ReceiveKeys(PublicKeys []*PublicKeyPair) {

	receiver.PublicKeys = make([]*PublicKeyPair, receiver.k)
	receiver.PublicKeys = PublicKeys
}

// Method to encrypt messages (seeds) when the parties invoke the κ×OTκ-functionality, where the OTSender plays the receiver and OTReceiver plays the sender.
func (receiver *OTReceiver) EncryptSeeds(elGamal *elgamal.ElGamal) []*CiphertextPair {

	k := receiver.k

	ciphertexts := make([]*CiphertextPair, k)

	for i := 0; i < k; i++ {

		ciphertexts[i] = &CiphertextPair{} // Initialize the ciphertext pair

		// Encrypt the messages using the public keys received from the OTReceiver
		ciphertexts[i].Ciphertext0 = elGamal.Encrypt(receiver.seeds[i].seed0, receiver.PublicKeys[i].MessageKey0)
		ciphertexts[i].Ciphertext1 = elGamal.Encrypt(receiver.seeds[i].seed1, receiver.PublicKeys[i].MessageKey1)
	}

	return ciphertexts

}

// Method for generating the bit matrix T of size m × κ, after the κ×OTκ OT-functionality, where the OTSender plays the receiver and OTReceiver plays the sender.
// GenerateMatrixT generates the bit matrix T after the k×OTk functionality.
func (receiver *OTReceiver) GenerateMatrixT() error {
	k := receiver.k
	m := receiver.m

	// Initialize the matrix T of size m × κ.
	T := make([][]string, m) // m rows.
	for i := range T {
		T[i] = make([]string, k) // k columns per row.
	}

	// Generate each column of T.
	for i := 0; i < k; i++ {
		// Generate a pseudo-random bitstring of m bits using the seed.
		bitstring, err := pseudoRandomGenerator(receiver.seeds[i].seed0, m)
		if err != nil {
			return err
		}
		print("Bitstring T: " + bitstring + "\n")

		for j := 0; j < m; j++ {
			T[j][i] = bitstring[j : j+1] // Assign the bit to the matrix T at position (j,i).
		}
	}

	print("T: " + "\n")
	PrintMatrix(T)

	// Assign the generated matrix to the receiver.
	receiver.T = T

	return nil
}

func (receiver *OTReceiver) GenerateAndSendMatrixU() [][]string {

	k := receiver.k
	m := receiver.m

	// Initialize the matrix U of size m × κ.
	U := make([][]byte, m) // m rows.
	for i := range U {
		U[i] = make([]byte, k) // k columns per row.
	}

	// Generate each column of U: u^i = t^i ⊕ G(k1_i ) ⊕ r.
	for i := 0; i < k; i++ {

		// Generate a pseudo-random bit string of m bits.
		bitstring, err := pseudoRandomGenerator(receiver.seeds[i].seed1, m)

		if err != nil {
			panic("Error from pseudoRandomGenerator in GenerateAndSendUMatrix: " + err.Error())
		}

		for j := 0; j < m; j++ {

			row_val, _ := strconv.Atoi(receiver.T[j][i])
			string_val, _ := strconv.Atoi(bitstring[j : j+1])
			U[j][i] = byte(row_val ^ string_val ^ receiver.selectionBits[i]) // u^i = t^i ⊕ G(k1_i ) ⊕ r = t^i ⊕ s_i ⊕ r

		}

	}

	return U

}

func (receiver *OTReceiver) DecryptMessage(ciphertextPairs []*CiphertextPair, elGamal *elgamal.ElGamal) []*big.Int {
	panic("implement me")
}

func (receiver *OTReceiver) ReceiveData(elGamal *elgamal.ElGamal) {
	panic("implement me")
}
