package handin5

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

	// Generate primes q and p such that p = kq + 1 for some k
	for {
		// Generate a large prime q of 256 bits length. Usually this should be around 2000 bits,
		// but for computation reasons we only use 256 bits.
		elGamal.q, _ = rand.Prime(rand.Reader, 256)

		elGamal.p = new(big.Int).Mul((big.NewInt(2)), elGamal.q) // p = kq (we use k = 2 for simplicity as said in the notes)
		elGamal.p = elGamal.p.Add(elGamal.p, big.NewInt(1))      // p = kq + 1

		if elGamal.p.ProbablyPrime(40) { // Test with 400 rounds of Miller-Rabin. Otherwise try new q and p values
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

func (elGamal *ElGamal) makeSecretKey() *big.Int {
	// sk ∈ [1, q-1]. Notice, we exclude 0 due to weak properties.
	qMinusOne := new(big.Int).Sub(elGamal.q, big.NewInt(1))
	x, _ := rand.Int(rand.Reader, qMinusOne) // random number x ∈ [0, q-2]
	sk := x.Add(x, big.NewInt(1))            // random sk ∈ [1, q-1]

	return sk
}

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

// The encrypt method first encodes the message m with the thrid encoding method from the notes:
// "check if (m + 1)^q = 1 mod p. If yes, encrypt M = m + 1. If not, encrypt M = −(m + 1)".
// Afterwards, it use the ElGamal encryption scheme to encrypt the encoded message M.
func (elGamal *ElGamal) Encrypt(m *big.Int, pk *big.Int) *Ciphertext {

	// Encoding of m
	var M *big.Int
	// Check if check if (m + 1)^q = 1 mod p.
	mPlusOne := new(big.Int).Add(m, big.NewInt(1))
	if new(big.Int).Exp(mPlusOne, elGamal.q, elGamal.p).Cmp(big.NewInt(1)) != 0 {
		M = mPlusOne // M = m + 1. If  (m + 1)^q = 1 mod p
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

	// Calculate M = c2 * (c1^sk)^-1 mod p
	s := new(big.Int).Exp(c1, sk, elGamal.p)        // s = c1^-sk mod p
	modInv := new(big.Int).ModInverse(s, elGamal.p) // modInv = s^-1 mod p
	M := new(big.Int).Mul(c2, modInv)               // M = c2 * s

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

// Exactly the same method as Init(), except that we use harcoded large primes for q which we found online.
// There seem to be a problem with the random number generator in GO, since the properties of the ElGamal cryptosystem
// does not hold for a prime found by rand.Prime(rand.Reader, 256).
func (elGamal *ElGamal) InitFixedQ() {

	// Hardcoded prime q found online at:
	// https://lists.exim.org/lurker/message/20200917.170121.9eb5c776.de.html
	// This is used due to primes found by rand.Prime(rand.Reader, 256) being buggy.
	qstr := "7FFFFFFFFFFFFFFFD6FC2A2C515DA54D57EE2B10139E9E78EC5CE2C1E7169B4AD4F09B208A3219FDE649CEE7124D9F7CBE97F1B1B1863AEC7B40D901576230BD69EF8F6AEAFEB2B09219FA8FAF83376842B1B2AA9EF68D79DAAB89AF3FABE49ACC278638707345BBF15344ED79F7F4390EF8AC509B56F39A98566527A41D3CBD5E0558C159927DB0E88454A5D96471FDDCB56D5BB06BFA340EA7A151EF1CA6FA572B76F3B1B95D8C8583D3E4770536B84F017E70E6FBF176601A0266941A17B0C8B97F4E74C2C1FFC7278919777940C1E1FF1D8DA637D6B99DDAFE5E17611002E2C778C1BE8B41D96379A51360D977FD4435A11C30942E4BFFFFFFFFFFFFFFFF"
	elGamal.q = new(big.Int)
	elGamal.q.SetString(qstr, 16)

	// Generate prime p such that p = kq + 1 for some k
	for {
		elGamal.p = new(big.Int).Mul((big.NewInt(2)), elGamal.q) // p = kq (we use k = 2 for simplicity as said in the notes)
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

// Test if the properties of the ElGamal cryptosystem holds for the public parameters p, q, g
// as explained in the cryptography course book by Ivan Damgård. The method was primarily used
// for debugging purposes with regards to the random number generator rand.Prime(rand.Reader, 256)
// as mentioned in the README. The methods prints an error if the properties does not hold.
func (elGamal *ElGamal) TestProperties() {
	// Check if q divides p-1
	pMinusOne := new(big.Int).Sub(elGamal.p, big.NewInt(1))
	if new(big.Int).Mod(pMinusOne, elGamal.q).Cmp(big.NewInt(0)) != 0 {
		fmt.Println("Error: q does not divide p-1")
	}

	// Check if g^q ≡ 1 (mod p)
	gToQ := new(big.Int).Exp(elGamal.g, elGamal.q, elGamal.p)
	if gToQ.Cmp(big.NewInt(1)) != 0 {
		fmt.Println("Error: g^q is not congruent to 1 mod p")
	}

	// Check if g^(q/2) is not congruent to 1 mod p
	qDiv2 := new(big.Int).Div(elGamal.q, big.NewInt(2))
	gToQDiv2 := new(big.Int).Exp(elGamal.g, qDiv2, elGamal.p)
	if gToQDiv2.Cmp(big.NewInt(1)) == 0 {
		fmt.Println("Error: g^(q/2) is congruent to 1 mod p ")
	}
}
