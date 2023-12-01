package elgamal

import (
	"crypto/rand"
	"math/big"
)

type ElGamal struct {
	q *big.Int // order of group G (cyclic subgroup of F_p). Notice, q | p-1
	p *big.Int // prime number defining the finite field F_p
	g *big.Int // generator of group G
}

// Ciphertext is a struct containing the two parts of a ciphertext (Due to GO being unable to return a list of tuple values).
// Used for encrypt since Bob needs to return a list of 2 ciphertexts each containing c1, c2.
type Ciphertext struct {
	C1 *big.Int
	C2 *big.Int
}

// Generate the public parameters p, q, g for the ElGamal cryptosystem
func (elGamal *ElGamal) Init() {

	// Generate primes q and p such that p = kq + 1 for some k
	for {
		// Generate a large prime q of 256 bits length. Usually this should be around 2048 bits ,
		// but for computation reasons we only use 256 bits.
		elGamal.q, _ = rand.Prime(rand.Reader, 2048)

		elGamal.p = new(big.Int).Mul((big.NewInt(2)), elGamal.q) // p = kq (we use k = 2 for simplicity as suggested in lecture notes)
		elGamal.p = elGamal.p.Add(elGamal.p, big.NewInt(1))      // p = kq + 1

		if elGamal.p.ProbablyPrime(400) { // Test with 400 rounds of Miller-Rabin. Otherwise try new q and p values
			break
		}

	}
	// Generate a DDH-safe group g of order q in Z_p^* by using the second suggesting from the notes:
	// "Pick arbitrary x from Z_p^* where x != 1 and x != -1, and compute g = x2 mod p".
	pMinusTwo := new(big.Int).Sub(elGamal.p, big.NewInt(1))   // pMinusTwo = p-2
	x, _ := rand.Int(rand.Reader, pMinusTwo)                  // random number x ∈ [0, p-3]
	x = x.Add(x, big.NewInt(2))                               // random number x ∈ [2, p-1]
	elGamal.g = new(big.Int).Exp(x, big.NewInt(2), elGamal.p) // g = x^2 mod p

}

func (elGamal *ElGamal) MakeSecretKey() *big.Int {
	// sk ∈ [1, q-1]. Notice, we exclude 0 due to weak properties.
	qMinusOne := new(big.Int).Sub(elGamal.q, big.NewInt(1))
	x, _ := rand.Int(rand.Reader, qMinusOne) // random number x ∈ [0, q-2]
	sk := x.Add(x, big.NewInt(1))            // random sk ∈ [1, q-1]

	return sk
}

// Generate a "real" public key h = g^sk mod p from a secret key sk
func (elGamal *ElGamal) Gen(sk *big.Int) *big.Int {

	h := new(big.Int).Exp(elGamal.g, sk, elGamal.p) // h = g^sk mod p

	return h // return public key
}

// OGen is the oblivious version of Gen. It returns a random "fake" public key following the second method in exercise 5
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

// The encrypt method first encodes the message m with the third encoding method from the notes:
// "check if (m + 1)^q = 1 mod p. If yes, encrypt M = m + 1. If not, encrypt M = −(m + 1)".
// Afterwards, it use the ElGamal encryption scheme to encrypt the encoded message M.
func (elGamal *ElGamal) Encrypt(m *big.Int, pk *big.Int) *Ciphertext {

	// Encoding of m
	var M *big.Int
	mPlusOne := new(big.Int).Add(m, big.NewInt(1))                                // m + 1
	if new(big.Int).Exp(mPlusOne, elGamal.q, elGamal.p).Cmp(big.NewInt(1)) == 0 { // Check If (m + 1)^q = 1 mod p
		M = mPlusOne // M = m + 1.
	} else {
		M = new(big.Int).Neg(mPlusOne) // M = -(m + 1)
	}

	// Generate a random number r ∈ [1, q-1]. Notice, we exclude 0 due to weak properties.
	qMinusOne := new(big.Int).Sub(elGamal.q, big.NewInt(1))
	r, _ := rand.Int(rand.Reader, qMinusOne) // random number ∈ [0, q-2]
	r = r.Add(r, big.NewInt(1))              // random number ∈ [1, q-1]

	c1 := new(big.Int).Exp(elGamal.g, r, elGamal.p)               // c1 = g^r mod p
	c2 := new(big.Int).Mul(M, new(big.Int).Exp(pk, r, elGamal.p)) // c2 = m * (pk^r mod p)
	c2 = c2.Mod(c2, elGamal.p)                                    // c2 = m * pk^r mod p

	return &Ciphertext{c1, c2} // return ciphertext struct since GO does not support tuple values

}

// Regular ElGamal decryption method with decoding afterwards. First M = c2 * (c1^sk)^-1 mod p is computed.
// Then the decoding method from the notes is used: "If M ≤ q, then m = M − 1, otherwise m = −M − 1."
func (elGamal *ElGamal) Decrypt(c1 *big.Int, c2 *big.Int, sk *big.Int) *big.Int {

	// Calculate M = c2 * c1^-sk mod p
	negSk := new(big.Int).Neg(sk)               // -sk
	s := new(big.Int).Exp(c1, negSk, elGamal.p) // s = c1^-sk mod p
	M := new(big.Int).Mul(c2, s)                // M = c2 * s
	M = M.Mod(M, elGamal.p)                     // M = c2 * s mod p

	// Decode M
	var m *big.Int
	// Check if M ≤ q and set m accordingly
	if M.Cmp(elGamal.q) <= 0 {
		m = new(big.Int).Sub(M, big.NewInt(1)) // m = M - 1. If M ≤ q
	} else {
		negatedM := new(big.Int).Neg(M)
		m = new(big.Int).Sub(negatedM, big.NewInt(1)) // m = -M - 1
	}

	return m.Mod(m, elGamal.p) // M = m mod p
}
