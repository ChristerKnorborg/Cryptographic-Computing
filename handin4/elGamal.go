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

	qstr := "7FFFFFFFFFFFFFFFD6FC2A2C515DA54D57EE2B10139E9E78EC5CE2C1E7169B4AD4F09B208A3219FDE649CEE7124D9F7CBE97F1B1B1863AEC7B40D901576230BD69EF8F6AEAFEB2B09219FA8FAF83376842B1B2AA9EF68D79DAAB89AF3FABE49ACC278638707345BBF15344ED79F7F4390EF8AC509B56F39A98566527A41D3CBD5E0558C159927DB0E88454A5D96471FDDCB56D5BB06BFA340EA7A151EF1CA6FA572B76F3B1B95D8C8583D3E4770536B84F017E70E6FBF176601A0266941A17B0C8B97F4E74C2C1FFC7278919777940C1E1FF1D8DA637D6B99DDAFE5E17611002E2C778C1BE8B41D96379A51360D977FD4435A11C30942E4BFFFFFFFFFFFFFFFF"
	pstr := "FFFFFFFFFFFFFFFFADF85458A2BB4A9AAFDC5620273D3CF1D8B9C583CE2D3695A9E13641146433FBCC939DCE249B3EF97D2FE363630C75D8F681B202AEC4617AD3DF1ED5D5FD65612433F51F5F066ED0856365553DED1AF3B557135E7F57C935984F0C70E0E68B77E2A689DAF3EFE8721DF158A136ADE73530ACCA4F483A797ABC0AB182B324FB61D108A94BB2C8E3FBB96ADAB760D7F4681D4F42A3DE394DF4AE56EDE76372BB190B07A7C8EE0A6D709E02FCE1CDF7E2ECC03404CD28342F619172FE9CE98583FF8E4F1232EEF28183C3FE3B1B4C6FAD733BB5FCBC2EC22005C58EF1837D1683B2C6F34A26C1B2EFFA886B423861285C97FFFFFFFFFFFFFFFF"
	gstr := "2"

	elGamal.q.SetString(qstr, 16)
	elGamal.p.SetString(pstr, 16)
	elGamal.g.SetString(gstr, 16)

	// fmt.Println("q: ", elGamal.q)
	// fmt.Println("p: ", elGamal.p)
	// fmt.Println("g: ", elGamal.g)

	// Check if properties hold for hardcoded values of p, q, g
	testProperties(elGamal.p, elGamal.q, elGamal.g)
}

func testProperties(p, q, g *big.Int) {
	// Check if q divides p-1
	pMinusOne := new(big.Int).Sub(p, big.NewInt(1))
	if new(big.Int).Mod(pMinusOne, q).Cmp(big.NewInt(0)) != 0 {
		fmt.Println("Error: q does not divide p-1")
	}

	// Check if g^q ≡ 1 (mod p)
	gToQ := new(big.Int).Exp(g, q, p)
	if gToQ.Cmp(big.NewInt(1)) != 0 {
		fmt.Println("Error: g^q is not congruent to 1 mod p")
	}

	// Check if g^(q/2) is not congruent to 1 mod p
	qDiv2 := new(big.Int).Div(q, big.NewInt(2))
	gToQDiv2 := new(big.Int).Exp(g, qDiv2, p)
	if gToQDiv2.Cmp(big.NewInt(1)) == 0 {
		fmt.Println("Error: g^(q/2) is congruent to 1 mod p ")
	}
}
