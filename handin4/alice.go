package handin4

import "math/big"

type Alice struct {
	q *big.Int // order of group G (cyclic subgroup of F_p). Notice, q | p-1
	p *big.Int // prime number defining the finite field F_p
	g *big.Int // generator of group G
}

func (alice *Alice) Init(p *big.Int, q *big.Int, g *big.Int) {
	alice.q = q
	alice.p = p
	alice.g = g
}

func (alice *Alice) MakePublicKeys(x int) {
	// Make list of 8 public keys
	publicKeys := make([]*big.Int, 8)

	pkList := make([]interface{}, 8)

	alice.x = x
}
