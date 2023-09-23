package handin4

import (
	"crypto/rand"
	"math/big"
)

type ElGamal struct {
	q *big.Int // order of group G (cyclic subgroup of F_p). Notice, q | p-1
	p *big.Int // prime number defining the finite field F_p
	g *big.Int // generator of group G

	r *big.Int // random number r ∈ [0, q-1]
}

// Ciphertext is a struct containing the two parts of a ciphertext (Due to GO being unable to return a list of tuple values).
// Used for encrypt since Bob needs to return a list of 8 ciphertexts each containing c1, c2.
type Ciphertext struct {
	c1 *big.Int
	c2 *big.Int
}

// Generate the public parameters p, q, g for the ElGamal cryptosystem
func (elGamal *ElGamal) Init() {
	elGamal.q, _ = rand.Prime(rand.Reader, 256) // Generate a large prime q of 256 bits length

	// Generate a prime p such that p = kq + 1 for some k
	for {
		k, _ := rand.Int(rand.Reader, big.NewInt(1<<16)) // choose a random k up to 2^16
		elGamal.p = new(big.Int).Mul(k, elGamal.q)
		elGamal.p = elGamal.p.Add(elGamal.p, big.NewInt(1))

		if elGamal.p.ProbablyPrime(40) { // Test with 40 rounds of Miller-Rabin. Otherwise repeat.
			break
		}
	}

	// Find a generator g of the subgroup of order q in Z_p^*
	for {
		elGamal.g, _ = rand.Int(rand.Reader, elGamal.p) // random number less than p
		if new(big.Int).Exp(elGamal.g, elGamal.q, elGamal.p).Cmp(big.NewInt(1)) == 0 && elGamal.g.Cmp(big.NewInt(1)) != 0 {
			break
		}
	}
}

func (elGamal *ElGamal) makeSecretKey() *big.Int {
	sk, _ := rand.Int(rand.Reader, elGamal.q) // sk ∈ [0, q-1]. Notice, we include 0 - even though it is technically a bad choice.
	return sk
}

func (elGamal *ElGamal) Gen(sk *big.Int) *big.Int {

	h := new(big.Int).Exp(elGamal.g, sk, elGamal.p) // h = g^sk mod p

	return h // return public key
}

// OGen is the oblivious version of Gen. It returns a random "fake" public key
func (elGamal *ElGamal) OGen() *big.Int {

	// Make a random big int r ∈ [0, p]
	r, _ := rand.Int(rand.Reader, elGamal.p)

	return r
}

func (elGamal *ElGamal) Encrypt(m *big.Int, pk *big.Int) *Ciphertext {
	// Generate a random number r ∈ [0, q-1]. Notice, we include 0 - even though it is technically a bad choice.
	r, _ := rand.Int(rand.Reader, elGamal.q)

	c1 := new(big.Int).Exp(elGamal.g, r, elGamal.p)               // c1 = g^r mod p
	c2 := new(big.Int).Mul(m, new(big.Int).Exp(pk, r, elGamal.p)) // c2 = m * pk^r mod p

	return &Ciphertext{c1, c2}

}

func (elGamal *ElGamal) Decrypt(c1 *big.Int, c2 *big.Int, sk *big.Int) *big.Int {
	// m = c2 * c1^(-sk) mod p

	s := new(big.Int).Exp(c1, sk, elGamal.p)      // s = c1^sk mod p
	sInv := new(big.Int).ModInverse(s, elGamal.p) // sInv = temp^-1 mod p.
	m := new(big.Int).Mul(c2, sInv)               // m = c2 * sInv
	return m.Mod(m, elGamal.p)                    // m = m mod p
}
