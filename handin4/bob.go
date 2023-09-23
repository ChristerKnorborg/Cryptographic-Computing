package handin4

import "math/big"

type Bob struct {
	q *big.Int // order of group G (cyclic subgroup of F_p). Notice, q | p-1
	p *big.Int // prime number defining the finite field F_p
	g *big.Int // generator of group G
}

func (bob *Bob) Init(p *big.Int, q *big.Int, g *big.Int) {
	bob.q = q
	bob.p = p
	bob.g = g
}
