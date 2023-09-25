package handin4

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

type ElGamal struct {
	q *big.Int // order of group G (cyclic subgroup of F_p). Notice, q | p-1
	p *big.Int // prime number defining the finite field F_p
	g *big.Int // generator of group G
}

// Ciphertext is a struct containing the two parts of a ciphertext (Due to GO being unable to return a list of tuple values).
// Used for encrypt since Bob needs to return a list of 8 ciphertexts each containing c1, c2.
type Ciphertext struct {
	c1 *big.Int
	c2 *big.Int
}

// Generate the public parameters p, q, g for the ElGamal cryptosystem
func (elGamal *ElGamal) Init() {

	// Generate a prime p such that p = kq + 1 for some k
	for {
		elGamal.q, _ = rand.Prime(rand.Reader, 256) // Generate a large prime q of 256 bits length
		// k, _ := rand.Int(rand.Reader, big.NewInt(1<<16))    // choose a random k up to 2^16 (2 for faster testing)
		elGamal.p = new(big.Int).Mul(big.NewInt(2), elGamal.q) // p = kq
		elGamal.p = elGamal.p.Add(elGamal.p, big.NewInt(1))    // p = kq + 1

		if elGamal.p.ProbablyPrime(400) { // Test with 400 rounds of Miller-Rabin. Otherwise repeat.
			break
		}
	}

	// Find a generator g of the subgroup of order q in Z_p^*
	for {
		// number should be in range [1, q-1]
		pMinusOne := new(big.Int).Sub(elGamal.p, big.NewInt(1)) // pMinusOne = p-1
		elGamal.g, _ = rand.Int(rand.Reader, pMinusOne)         // random number ∈ [0, p-2]
		elGamal.g = elGamal.g.Add(elGamal.g, big.NewInt(1))     // random number ∈ [1, p-1]

		// Condition 1: g != 1
		if elGamal.g.Cmp(big.NewInt(1)) == 0 {
			continue // try again
		}

		// Condition 2: g^q mod p = 1
		if new(big.Int).Exp(elGamal.g, elGamal.q, elGamal.p).Cmp(big.NewInt(1)) != 0 {
			continue // try again
		}

		// Condition 3: g^(p-1)/q mod p != 1
		order := new(big.Int).Div(elGamal.p.Sub(elGamal.p, big.NewInt(1)), elGamal.q)
		if new(big.Int).Exp(elGamal.g, order, elGamal.p).Cmp(big.NewInt(1)) == 0 {
			continue // try again
		}

		break
	}

}

func (elGamal *ElGamal) makeSecretKey() *big.Int {
	// sk ∈ [1, q-1]. Notice, we exclude 0 due to weak properties.
	qMinusOne := new(big.Int).Sub(elGamal.q, big.NewInt(1))
	sk, _ := rand.Int(rand.Reader, qMinusOne) // random number ∈ [0, q-2]
	sk = sk.Add(sk, big.NewInt(1))            // random number ∈ [1, q-1]

	return sk
}

func (elGamal *ElGamal) Gen(sk *big.Int) *big.Int {

	h := new(big.Int).Exp(elGamal.g, sk, elGamal.p) // h = g^sk mod p

	return h // return public key
}

// OGen is the oblivious version of Gen. It returns a random "fake" public key following the 2. exercise 5 point.
func (elGamal *ElGamal) OGen() *big.Int {

	n := elGamal.p.BitLen() // Get the bit length of p

	// Create 2^(2n) upper bound for random number
	upperBound := new(big.Int).Lsh(big.NewInt(1), uint(2*n)) // Left shift is equivalent to multiplying by 2 raised to a power

	// Generate a random big int r ∈ [1, 2^(2n)]
	r, _ := rand.Int(rand.Reader, upperBound)
	if r.Sign() == 0 {
		r = r.Add(r, big.NewInt(1)) // Ensure r is not zero
	}

	// Return r mod p
	return r.Mod(r, elGamal.p)
}
func (elGamal *ElGamal) Encrypt(m *big.Int, pk *big.Int) *Ciphertext {
	// Generate a random number r ∈ [1, q-1]. Notice, we exclude 0 due to weak properties.
	qMinusOne := new(big.Int).Sub(elGamal.q, big.NewInt(1))
	r, _ := rand.Int(rand.Reader, qMinusOne) // random number ∈ [0, q-2]
	r = r.Add(r, big.NewInt(1))              // random number ∈ [1, q-1]

	c1 := new(big.Int).Exp(elGamal.g, r, elGamal.p)               // c1 = g^r mod p
	c2 := new(big.Int).Mul(m, new(big.Int).Exp(pk, r, elGamal.p)) // c2 = m * (pk^r mod p)
	c2 = c2.Mod(c2, elGamal.p)                                    // c2 = m * pk^r mod p

	return &Ciphertext{c1, c2}

}

// m = c2 * (c1^sk)^-1 mod p
func (elGamal *ElGamal) Decrypt(c1 *big.Int, c2 *big.Int, sk *big.Int) *big.Int {

	s := new(big.Int).Exp(c1, sk, elGamal.p)        // s = c1^-sk mod p
	modInv := new(big.Int).ModInverse(s, elGamal.p) // modInv = s^-1 mod p

	m := new(big.Int).Mul(c2, modInv) // m = c2 * s
	return m.Mod(m, elGamal.p)        // m = m mod p
}

func (elGamal *ElGamal) InitFixedValues() {
	// Fix the public parameters p, q, g for the ElGamal cryptosystem
	elGamal.q = new(big.Int)
	elGamal.p = new(big.Int)
	elGamal.g = new(big.Int)

	elGamal.q.SetString("FCA682CE8E12CABA26EFCCF7110E526DB078B05EDECBCD1EB4A208F3AE1617AE01F", 16)
	elGamal.p.SetString("9A4D6792295A090D5771133B60E8AB2FC6C4E1299B6FE6C3B6E46BDE4107A4C954FEB4FD9ADA8E6DD04F80B3EF2D9ED4B304F564A1C895C57B49509B758AC841B6E65B877C3109D27A7AA0E6A42E294F0E4E507663DB5F504CF3B80887FAA7BE8EF590FE7B3E8CE4A7B4D69C15B5FCB642", 16)
	elGamal.g.SetString("2", 16)

	fmt.Println("q: ", elGamal.q)
	fmt.Println("p: ", elGamal.p)
	fmt.Println("g: ", elGamal.g)

}
