package handin6

import (
	"crypto/rand"
	"math/big"
)

// Depth d Homomorphic Encryption (DHE) scheme
type DHE struct {
	p *big.Int   // secret key
	q []*big.Int // random integers q_1,..., q_n
	r []*big.Int // random small integers r_1,..., r_n.
	y []*big.Int // public key (y1,..., yn) where y_i = p * q_i + 2r_i
	n int
}

func (dhe *DHE) Init() {

}

func (dhe *DHE) GenerateKeys() {

	// Choose a random big integer of length n. Notice, left shift is equivalent to multiplying by 2 raised to a power
	p, _ := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), uint(dhe.n)))

	p.Or(p, big.NewInt(1)) // Ensure p is odd by setting the least significant bit to 1
	dhe.p = p

	for i := 0; i < dhe.n; i++ {
		// Choose random big q value of length n
		q_i, _ := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), uint(dhe.n)))
		dhe.q[i] = q_i

		// Choose small random r value of length n/10
		smallNumber := new(big.Int).Lsh(big.NewInt(1), uint(dhe.n/10)) // 10% of n's bit length
		r_i, _ := rand.Int(rand.Reader, smallNumber)
		dhe.r[i] = r_i

		// Calculate y_i = p * q_i + 2r_i
		pMulQ_i := new(big.Int).Mul(dhe.p, q_i)        // p * q_i
		twoR_i := new(big.Int).Mul(r_i, big.NewInt(2)) // 2r_i

		y_i := new(big.Int).Add(pMulQ_i, twoR_i) // y_i = p * q_i + 2r_i
		dhe.y[i] = y_i
	}

}

func (dhe *DHE) Encrypt(m int) *big.Int {

	// Check if m is a bit (0 or 1)
	if m != 0 && m != 1 {
		panic("m must be a bit (0 or 1)")
	}

	// Create a random subset S of {1, ..., n}
	S := randomSubset(dhe.n)

	// Initialize c with the value of m
	mInt64 := int64(m) // Convert m to int64 (This is due to GO BigInt compatibility for NewInt)
	c := big.NewInt(mInt64)

	// Compute c = m + ∑_(i∈S) y_i
	for _, i := range S {
		c.Add(c, dhe.y[i])
	}

	return c
}

func (dhe *DHE) Decrypt(c *big.Int) int {
	// Compute c mod p
	cModP := new(big.Int).Mod(c, dhe.p)

	// Compute (c mod p) mod 2
	mInt64 := new(big.Int).Mod(cModP, big.NewInt(2))

	// Convert m to int from int64 (Int64 in the first place is due to GO BigInt compatibility for Int64)
	m := int(mInt64.Int64())
	return m
}

// This function creates a random subset of {1,..., n} by iterating through each number
// from 0 to n−1 and deciding with a 50% chance whether to include that number in the subset or not.
func randomSubset(n int) []int {
	subset := make([]int, 0)
	for i := 1; i <= n; i++ {
		// Decide with 50% probability whether to include i in the subset or not
		randomBit, _ := rand.Int(rand.Reader, big.NewInt(2)) // randomBit ∈ {0, 1}
		if randomBit.Int64() == 1 {
			subset = append(subset, i)
		}
	}
	return subset
}

// Homomorphic function applies XOR between two ciphertexts c1 and c2 by addition and returns the result ciphertext
func HE_xor(c1 *big.Int, c2 *big.Int) *big.Int {
	return new(big.Int).Add(c1, c2)
}

// Homomorphic function applies AND between two ciphertexts c1 and c2 by multiplication and returns the result ciphertext
func HE_and(c1 *big.Int, c2 *big.Int) *big.Int {
	return new(big.Int).Mul(c1, c2)
}
