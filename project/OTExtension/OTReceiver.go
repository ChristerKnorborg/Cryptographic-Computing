// OTReceiver.go
package OTExtension

// Import your ElGamal package
import (
	"crypto/rand"
	"cryptographic-computing/project/elgamal"
	"math/big"
	"strconv"

	"github.com/hashicorp/vault/sdk/helper/xor"
)

type OTReceiver struct {
	m             int              // Number of messages to be received
	k             int              // Security parameter
	l             int              // Bit length of each message
	SelectionBits []int            // Receiver R holds m selection bits r = (r_1, ..., r_m). AFTER DEBUG CHANGE TO SMALL s.
	Seeds         []*Seed          // Messages (Seeds) to be sent, when invoking the κ×OTκ-functionality, where the OTSender plays the receiver and OTReceiver plays the sender. //CHANGE AFTER DEBUG
	secretKeys    []*big.Int       // Secret keys for each message to be received.
	PublicKeys    []*PublicKeyPair // Public keys received from the OTSender when invoking the κ×OTκ-functionality, where the OTSender plays the receiver and OTReceiver plays the sender.
	T             [][]string       // Bit matrix T of size m × κ, after the κ×OTκ OT-functionality
}

func (receiver *OTReceiver) Init(selectionBits []int, securityParameter int, l int) {

	receiver.l = l
	receiver.m = len(selectionBits)
	receiver.SelectionBits = selectionBits
	receiver.k = securityParameter

}

// The receiver chooses k pairs of k-bit Seeds {(k0_i , k1_i )} from i = 1 to k
func (receiver *OTReceiver) ChooseSeeds() {

	k := receiver.k

	Seeds := make([]*Seed, k)

	for i := 0; i < k; i++ {
		// Generate a k-bit random number for Seed0
		Seed0, err := rand.Int(rand.Reader, big.NewInt(1).Lsh(big.NewInt(1), uint(k)))
		if err != nil {
			panic("Error in ChooseSeeds for Seed0: " + err.Error())
		}
		// Generate a k-bit random number for Seed1
		Seed1, err := rand.Int(rand.Reader, big.NewInt(1).Lsh(big.NewInt(1), uint(k)))
		if err != nil {
			panic("Error in ChooseSeeds for Seed1: " + err.Error())
		}
		Seeds[i] = &Seed{
			Seed0: Seed0,
			Seed1: Seed1,
		}
	}
	receiver.Seeds = Seeds

	// print("printing the Seeds in ChooseSeeds\n")
	// for _, s := range Seeds {
	// 	print("Seed0", s.Seed0.String()+"\n")
	// 	print("Seed1", s.Seed1.String()+"\n")

	// }
}

// Method for to receive Public keys, when the parties invoke the κ×OTκ-functionality, where the OTSender plays the receiver and OTReceiver plays the sender.
func (receiver *OTReceiver) ReceiveKeys(PublicKeys []*PublicKeyPair) {

	receiver.PublicKeys = make([]*PublicKeyPair, receiver.k)
	receiver.PublicKeys = PublicKeys
}

// Method to encrypt messages (Seeds) when the parties invoke the κ×OTκ-functionality, where the OTSender plays the receiver and OTReceiver plays the sender.
func (receiver *OTReceiver) EncryptSeeds(elGamal *elgamal.ElGamal) []*CiphertextPair {

	k := receiver.k

	ciphertextPairs := make([]*CiphertextPair, k)

	for i := 0; i < k; i++ {

		ciphertextPairs[i] = &CiphertextPair{} // Initialize the ciphertext pair

		// Encrypt the messages using the public keys received from the OTSender(receiver in initial phase)
		ciphertextPairs[i].Ciphertext0 = elGamal.Encrypt(receiver.Seeds[i].Seed0, receiver.PublicKeys[i].MessageKey0)
		ciphertextPairs[i].Ciphertext1 = elGamal.Encrypt(receiver.Seeds[i].Seed1, receiver.PublicKeys[i].MessageKey1)
	}

	return ciphertextPairs

}

// Method for generating the bit matrix T of size m × κ, after the κ×OTκ OT-functionality, where the OTSender plays the receiver and OTReceiver plays the sender.
// GenerateMatrixT generates the bit matrix T after the k×OTk functionality.
func (receiver *OTReceiver) GenerateMatrixT() [][]string {
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
		bitstring, err := pseudoRandomGenerator(receiver.Seeds[i].Seed0, m)
		print("bitstring in T: " + bitstring + "\n")

		if err != nil {
			panic("Error from pseudoRandomGenerator in GenerateMatrixT: " + err.Error())
		}

		for j := 0; j < m; j++ {
			T[j][i] = bitstring[j : j+1] // Assign the bit to the matrix T at position (j,i).
		}
	}
	// Assign the generated matrix to the receiver.
	receiver.T = T

	print("Matrix T: \n")
	PrintMatrix(T)
	return T
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
		bitstring, err := pseudoRandomGenerator(receiver.Seeds[i].Seed1, m)
		print("bitstring in U: " + bitstring + "\n")

		if err != nil {
			panic("Error from pseudoRandomGenerator in GenerateAndSendUMatrix: " + err.Error())
		}

		for j := 0; j < m; j++ {

			T_idx := receiver.T[j][i]
			G_idx := bitstring[j : j+1]
			selection_bit := receiver.SelectionBits[j]

			xor, err := XOR(T_idx, G_idx, selection_bit)
			if err != nil {
				panic("Error from XOR in GenerateAndSendUMatrix: " + err.Error())
			}
			U[j][i] = xor

			print("T_idx, G_idx, sel_bit " + " " + T_idx + " " + G_idx + " " + strconv.Itoa(selection_bit) + "\n")
			print("xor " + xor + "\n")
		}
	}
	print("Matrix U: \n")
	PrintMatrix(U)
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
		if receiver.SelectionBits[j] == 0 {
			y_j = ByteCiphertextPairs[j].y0
		} else if receiver.SelectionBits[j] == 1 {
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
		print("hash receiver : \n")
		PrintBinaryString(hash)

		xor, err := xor.XORBytes(y_j, hash)
		if err != nil {
			panic("Error from XOR in DecryptCiphertexts: " + err.Error())
		}
		plaintexts[j] = xor

	}
	// Print the result
	print("printing the result in DecryptCiphertexts\n")
	for _, b := range plaintexts {
		PrintBinaryString(b)
	}
	return plaintexts
}
