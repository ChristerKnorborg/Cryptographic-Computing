// OTReceiver.go
package OTExtension

// Import your ElGamal package
import (
	"crypto/rand"
	"cryptographic-computing/project/elgamal"
	"cryptographic-computing/project/utils"
	"math/big"

	"github.com/hashicorp/vault/sdk/helper/xor"
)

type OTReceiver struct {
	m             int                    // Number of messages to be received
	k             int                    // Security parameter
	l             int                    // Bit length of each message
	selectionBits []int                  // Receiver R holds m selection bits r = (r_1, ..., r_m).
	seeds         []*utils.Seed          // Messages (seeds) to be sent, when invoking the κ×OTκ-functionality, where the OTSender plays the receiver and OTReceiver plays the sender.
	PublicKeys    []*utils.PublicKeyPair // Public keys received from the OTSender when invoking the κ×OTκ-functionality, where the OTSender plays the receiver and OTReceiver plays the sender.
	T             [][]uint8              // Bit matrix T of size m × κ, after the κ×OTκ OT-functionality
}

func (receiver *OTReceiver) Init(selectionBits []int, securityParameter int, l int) {

	receiver.l = l
	receiver.m = len(selectionBits)
	receiver.selectionBits = selectionBits
	receiver.k = securityParameter

}

// The receiver chooses k pairs of k-bit seeds {(k0_i , k1_i )} from i = 1 to k
func (receiver *OTReceiver) ChooseSeeds() {

	k := receiver.k

	seeds := make([]*utils.Seed, k)

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

		seeds[i] = &utils.Seed{
			Seed0: seed0,
			Seed1: seed1,
		}
	}

	receiver.seeds = seeds

}

// Method for to receive Public keys, when the parties invoke the κ×OTκ-functionality, where the OTSender plays the receiver and OTReceiver plays the sender.
func (receiver *OTReceiver) ReceiveKeys(PublicKeys []*utils.PublicKeyPair) {

	receiver.PublicKeys = make([]*utils.PublicKeyPair, receiver.k)
	receiver.PublicKeys = PublicKeys
}

// Method to encrypt messages (seeds) when the parties invoke the κ×OTκ-functionality, where the OTSender plays the receiver and OTReceiver plays the sender.
func (receiver *OTReceiver) EncryptSeeds(elGamal *elgamal.ElGamal) []*utils.CiphertextPair {

	k := receiver.k

	ciphertextPairs := make([]*utils.CiphertextPair, k)

	for i := 0; i < k; i++ {

		ciphertextPairs[i] = &utils.CiphertextPair{} // Initialize the ciphertext pair

		// Encrypt the messages using the public keys received from the OTSender(receiver in initial phase)
		ciphertextPairs[i].Ciphertext0 = elGamal.Encrypt(receiver.seeds[i].Seed0, receiver.PublicKeys[i].MessageKey0)
		ciphertextPairs[i].Ciphertext1 = elGamal.Encrypt(receiver.seeds[i].Seed1, receiver.PublicKeys[i].MessageKey1)
	}

	return ciphertextPairs

}

// Method for generating the bit matrix T of size m × κ, after the κ×OTκ OT-functionality, where the OTSender plays the receiver and OTReceiver plays the sender.
// GenerateMatrixT generates the bit matrix T after the k×OTk functionality.
func (receiver *OTReceiver) GenerateMatrixT() {
	k := receiver.k
	m := receiver.m

	// Initialize the matrix T of size m × κ.
	T := make([][]uint8, m) // m rows.
	for i := range T {
		T[i] = make([]uint8, k) // k columns per row.
	}

	// Generate each column of T.
	for i := 0; i < k; i++ {
		// Generate a pseudo-random bitstring of m bits using the seed.
		bitstring, err := utils.PseudoRandomGenerator(receiver.seeds[i].Seed0, m)

		if err != nil {
			panic("Error from pseudoRandomGenerator in GenerateMatrixT: " + err.Error())
		}

		for j := 0; j < m; j++ {
			T[j][i] = bitstring[j] // Assign the bit to the matrix T at position (j,i).
		}
	}
	// Assign the generated matrix to the receiver.
	receiver.T = T

}

func (receiver *OTReceiver) GenerateAndSendMatrixU() [][]uint8 {

	k := receiver.k
	m := receiver.m

	// Initialize the matrix U of size m × κ.
	U := make([][]uint8, m) // m rows.
	for i := range U {
		U[i] = make([]uint8, k) // k columns per row.
	}

	// Generate each column of U: u^i = t^i ⊕ G(k1_i ) ⊕ r.
	for i := 0; i < k; i++ {

		// Generate a pseudo-random bit string of m bits.
		bitstring, err := utils.PseudoRandomGenerator(receiver.seeds[i].Seed1, m)

		if err != nil {
			panic("Error from pseudoRandomGenerator in GenerateAndSendUMatrix: " + err.Error())
		}

		for j := 0; j < m; j++ {

			T_idx := receiver.T[j][i]
			G_idx := bitstring[j]
			selection_bit := receiver.selectionBits[j]

			U[j][i] = T_idx ^ G_idx ^ uint8(selection_bit)
		}
	}
	return U
}

// The receiver then computes x^(r_j)_j = y^(rj)_j ⊕ H(j, t_j) for every 1 ≤ j ≤ m.
func (receiver *OTReceiver) DecryptCiphertexts(ByteCiphertextPairs []*utils.ByteCiphertextPair) [][]byte {

	m := receiver.m
	k := receiver.k
	l := receiver.l

	// Initialize the result
	plaintexts := make([][]byte, m)

	for j := 0; j < m; j++ {

		var y_j []byte
		if receiver.selectionBits[j] == 0 {
			y_j = ByteCiphertextPairs[j].Y0
		} else if receiver.selectionBits[j] == 1 {
			y_j = ByteCiphertextPairs[j].Y1
		} else {
			panic("Receiver choice bits are not 0 or 1 in DecryptCiphertexts")
		}

		t_row := make([]uint8, k)
		for i := 0; i < k; i++ {
			t_row[i] = receiver.T[j][i]
		}
		hash := utils.Hash(t_row, l)

		xor, err := xor.XORBytes(y_j, hash)
		if err != nil {
			panic("Error from XOR in DecryptCiphertexts: " + err.Error())
		}
		plaintexts[j] = xor

	}
	return plaintexts
}
