// OTReceiver.go
package OTExtension

// Import your ElGamal package
import (
	"crypto/rand"
	"cryptographic-computing/project/elgamal"
	"fmt"
	"math/big"

	"github.com/hashicorp/vault/sdk/helper/xor"
)

type OTReceiver struct {
	m             int              // Number of messages to be received
	k             int              // Security parameter
	l             int              // Bit length of each message
	selectionBits []int            // Receiver R holds m selection bits r = (r_1, ..., r_m).
	seeds         []*Seed          // Messages (seeds) to be sent, when invoking the κ×OTκ-functionality, where the OTSender plays the receiver and OTReceiver plays the sender.
	secretKeys    []*big.Int       // Secret keys for each message to be received.
	PublicKeys    []*PublicKeyPair // Public keys received from the OTSender when invoking the κ×OTκ-functionality, where the OTSender plays the receiver and OTReceiver plays the sender.
	T             [][]string       // Bit matrix T of size m × κ, after the κ×OTκ OT-functionality
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

	seeds := make([]*Seed, k)

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
		seeds[i] = &Seed{
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

		for j := 0; j < m; j++ {
			T[j][i] = bitstring[j : j+1] // Assign the bit to the matrix T at position (j,i).
		}
	}
	// Assign the generated matrix to the receiver.
	receiver.T = T

	return nil
}

func (receiver *OTReceiver) GenerateAndSendMatrixU() [][]string {

	k := receiver.k
	m := receiver.m

	// Initialize the matrix U of size m × κ.
	U := make([][]string, m) // m rows.
	for i := range U {
		U[i] = make([]string, k) // k columns per row.
	}

	// Generate each column of U: u^i = t^i ⊕ G(k1_i ) ⊕ r.
	for i := 0; i < k; i++ {

		// Generate a pseudo-random bit string of m bits.
		bitstring, err := pseudoRandomGenerator(receiver.seeds[i].seed1, m)

		if err != nil {
			panic("Error from pseudoRandomGenerator in GenerateAndSendUMatrix: " + err.Error())
		}

		for j := 0; j < m; j++ {

			T_idx := receiver.T[j][i]
			Q_idx := bitstring[j : j+1]
			selection_bit := receiver.selectionBits[i]

			xor, err := XOR(T_idx, Q_idx, selection_bit)
			if err != nil {
				panic("Error from XOR in GenerateAndSendUMatrix: " + err.Error())
			}
			U[j][i] = xor
		}
	}

	return U

}

// The receiver then computes x^(r_j)_j = y^(rj)_j ⊕ H(j, t_j) for every 1 ≤ j ≤ m.
func (receiver *OTReceiver) DecryptCiphertexts(ByteCiphertextPairs []*ByteCiphertextPair) [][]byte {

	m := receiver.m
	k := receiver.k
	l := receiver.l

	// Initialize the result
	plaintexts := make([][]byte, m)

	for j := 0; j < m; j++ {

		var y_j []byte
		if receiver.selectionBits[j] == 0 {
			y_j = ByteCiphertextPairs[j].y0
		} else if receiver.selectionBits[j] == 1 {
			y_j = ByteCiphertextPairs[j].y1
		} else {
			panic("Receiver choice bits are not 0 or 1 in DecryptCiphertexts")
		}

		t_row := ""
		for i := 0; i < k; i++ {
			t_idx := receiver.T[j][i]
			t_row += t_idx
		}
		hash := Hash([]byte(t_row), l)
		print("hash receiver " + string(hash) + "\n")

		xor, err := xor.XORBytes(y_j, hash)
		if err != nil {
			panic("Error from XOR in DecryptCiphertexts: " + err.Error())
		}
		plaintexts[j] = xor

	}
	// Print the result
	print("printing the result in DecryptCiphertexts\n")
	for _, b := range plaintexts {
		fmt.Printf("%08b ", b) // Binary
	}
	return plaintexts
}
