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
	m             int                    // Number of messages to be received.
	k             int                    // Security parameter.
	l             int                    // Bit length of each message.
	selectionBits []byte                 // Receiver R holds m selection bits r = (r_1, ..., r_m).
	seeds         []*utils.Seed          // Messages (seeds) to be sent, when invoking the Regular OT functionality k times.
	PublicKeys    []*utils.PublicKeyPair // Public keys received from the OTSender when invoking the Regular OT functionality k times.
	T             [][]byte               // Bit matrix T of size m × κ.
}

func (receiver *OTReceiver) Init(selectionBits []byte, k int, l int) {

	receiver.l = l
	receiver.m = len(selectionBits)
	receiver.selectionBits = selectionBits
	receiver.k = k
}

// The receiver chooses k pairs of k-bit seeds {(k0_i , k1_i )} from i = 1 to k using a secure random number generator.
func (receiver *OTReceiver) ChooseSeeds() {

	k := receiver.k

	seeds := make([]*utils.Seed, receiver.k)

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

// DEBUGGING METHOD. The receiver chooses k pairs of k-bit seeds {(k0_i , k1_i )} from i = 1 to k using fixed values.
func (receiver *OTReceiver) ChooseFixedSeeds() {

	k := receiver.k

	// Define fixed seeds
	fixedSeeds := []struct {
		Seed0 *big.Int
		Seed1 *big.Int
	}{
		// Fixed values for seeds
		{big.NewInt(12345), big.NewInt(67890)},
	}

	seeds := make([]*utils.Seed, k)

	for i := 0; i < k; i++ {
		// Use fixed seeds instead of generating them and ake sure the number of fixed seeds is at least k
		seed0 := fixedSeeds[i%len(fixedSeeds)].Seed0
		seed1 := fixedSeeds[i%len(fixedSeeds)].Seed1

		seeds[i] = &utils.Seed{
			Seed0: seed0,
			Seed1: seed1,
		}
	}
	receiver.seeds = seeds
}

// Method to receive Public keys, when the parties invoke the regular OT functionality k times,
// where the OTSender plays the receiver and OTReceiver plays the sender.
func (receiver *OTReceiver) ReceiveKeys(PublicKeys []*utils.PublicKeyPair) {

	receiver.PublicKeys = make([]*utils.PublicKeyPair, receiver.k)
	receiver.PublicKeys = PublicKeys
}

// Method to encrypt messages (seeds) when the parties invoke the regular OT functionality k times,
// where the OTSender plays the receiver and OTReceiver plays the sender.
func (receiver *OTReceiver) EncryptSeeds(elGamal *elgamal.ElGamal) []*utils.CiphertextPair {

	k := receiver.k

	ciphertextPairs := make([]*utils.CiphertextPair, k)

	for i := 0; i < k; i++ {

		ciphertextPairs[i] = &utils.CiphertextPair{} // Initialize the ciphertext pair

		// Encrypt the messages using the public keys received from the OTSender(who is receiver in this phase)
		ciphertextPairs[i].Ciphertext0 = elGamal.Encrypt(receiver.seeds[i].Seed0, receiver.PublicKeys[i].MessageKey0)
		ciphertextPairs[i].Ciphertext1 = elGamal.Encrypt(receiver.seeds[i].Seed1, receiver.PublicKeys[i].MessageKey1)
	}
	return ciphertextPairs
}

// Method for generating the bit matrices T and U of size m × κ.
// Notice this method is inefficient, since it accesses the entire matrix T (e.g. m x k entries).
func (receiver *OTReceiver) GenerateMatrixTAndU() [][]byte {
	k := receiver.k
	m := receiver.m

	// Initialize the matrix T of size m × κ.
	T := make([][]byte, m) // m rows.
	U := make([][]byte, m) // m rows.
	for i := range T {
		T[i] = make([]byte, k) // k columns per row.
		U[i] = make([]byte, k) // k columns per row.
	}

	// Generate each column of T.
	for i := 0; i < k; i++ {
		// Generate pseudo-random bitstrings of m bits using the seeds
		bitstringT, err1 := utils.PseudoRandomGenerator(receiver.seeds[i].Seed0, m)
		bitstringU, err2 := utils.PseudoRandomGenerator(receiver.seeds[i].Seed1, m)

		if err1 != nil || err2 != nil {
			panic("Error from pseudoRandomGenerator in GenerateMatrixT: " + err1.Error() + err2.Error())
		}

		for j := 0; j < m; j++ {
			T[j][i] = bitstringT[j] // Assign the bit to the matrix T at position (j,i).

			T_idx := T[j][i]
			G_idx := bitstringU[j]
			selection_bit := receiver.selectionBits[j]
			U[j][i] = T_idx ^ G_idx ^ selection_bit
		}
	}
	// Assign the generated matrix to the receiver.
	receiver.T = T
	return U // Send U to the OTSender
}

// More efficient method than the previous one for generating the bit matrices T and U of size m × κ.
// It generates the matrix T and U row-wise for transposing afterwards.
// Notice, only T is transposed here, as U is needed for the OTSender to generate Q row-wise.
func (receiver *OTReceiver) GenerateMatrixTAndUTranspose() [][]byte {
	k := receiver.k
	m := receiver.m

	// Initialize the matrix T and U of size k × m (transposed later)
	T := make([][]byte, k) // m rows.
	U := make([][]byte, k) // m rows.
	for i := range T {
		T[i] = make([]byte, m) // k columns per row.
		U[i] = make([]byte, m) // k columns per row.
	}

	for i := 0; i < k; i++ {
		// Generate pseudo-random bitstrings of m bits using the seeds
		bitstringT, err1 := utils.PseudoRandomGenerator(receiver.seeds[i].Seed0, m)
		bitstringU, err2 := utils.PseudoRandomGenerator(receiver.seeds[i].Seed1, m)
		if err1 != nil || err2 != nil {
			panic("Error from pseudoRandomGenerator in GenerateMatrixTAndUTranspose: " + err1.Error() + err2.Error())
		}
		T[i] = bitstringT

		xor1, err1 := xor.XORBytes(T[i], bitstringU)
		xor2, err2 := xor.XORBytes(xor1, receiver.selectionBits)
		if err1 != nil || err2 != nil {
			panic("Error from XOR in GenerateMatrixTAndUTranspose: " + err1.Error() + err2.Error())
		}
		U[i] = xor2
	}
	// Assign the generated matrix to the receiver after transposition.
	receiver.T = utils.TransposeMatrix(T)

	return U // Send U to the OTSender
}

// Even more efficient Method than the previous one for generating the bit matrices T and U of size m × κ.
// It generates the matrix T and U row-wise for transposing afterwards using Eklundh's algorithm.
// Notice, only T is transposed here, as U is needed for the OTSender to generate Q row-wise.
func (receiver *OTReceiver) GenerateMatrixTAndUEklundh(multithreaded bool) [][]byte {
	k := receiver.k
	m := receiver.m

	// Initialize the matrix T and U of size k × m (transposed later).
	T := make([][]byte, k) // m rows.
	U := make([][]byte, k) // m rows.
	for i := range T {
		T[i] = make([]byte, m) // k columns per row.
		U[i] = make([]byte, m) // k columns per row.
	}

	for i := 0; i < k; i++ {
		// Generate pseudo-random bitstrings of m bits using the seed.
		bitstringT, err1 := utils.PseudoRandomGenerator(receiver.seeds[i].Seed0, m)
		bitstringU, err2 := utils.PseudoRandomGenerator(receiver.seeds[i].Seed1, m)
		if err1 != nil || err2 != nil {
			panic("Error from pseudoRandomGenerator in GenerateMatrixTEklundh: " + err1.Error() + err2.Error())
		}
		T[i] = bitstringT

		xor1, err1 := xor.XORBytes(T[i], bitstringU)
		xor2, err2 := xor.XORBytes(xor1, receiver.selectionBits)
		if err1 != nil || err2 != nil {
			panic("Error from XOR in GenerateMatrixTAndUEklundh: " + err1.Error() + err2.Error())
		}
		U[i] = xor2
	}
	// Assign the generated matrix to the receiver, where Eklundh's algorithm is used to transpose .
	receiver.T = utils.EklundhTranspose(T, multithreaded)

	return U // Send U to the OTSender
}

// Method for decrypting the ciphertexts received from the OTSender.
// The receiver computes x^(r_j)_j = y^(r_j)_j ⊕ H(j, t_j) for every 1 ≤ j ≤ m.
func (receiver *OTReceiver) DecryptCiphertexts(ByteCiphertextPairs []*utils.ByteCiphertextPair) [][]byte {

	m := receiver.m
	l := receiver.l

	plaintexts := make([][]byte, m)

	for j := 0; j < m; j++ {

		var y_j []byte

		// Choose the ciphertext to decrypt based on the selection bit.
		if receiver.selectionBits[j] == 0 {
			y_j = ByteCiphertextPairs[j].Y0
		} else if receiver.selectionBits[j] == 1 {
			y_j = ByteCiphertextPairs[j].Y1
		} else {
			panic("Receiver choice bits are not 0 or 1 in DecryptCiphertexts")
		}

		hash := utils.Hash(receiver.T[j], l) // Generate hash of length l from the j'th row of T.

		xor, err := xor.XORBytes(y_j, hash) // XOR the ciphertext with the hash.
		if err != nil {
			panic("Error from XOR in DecryptCiphertexts: " + err.Error())
		}
		plaintexts[j] = xor

	}
	return plaintexts
}
