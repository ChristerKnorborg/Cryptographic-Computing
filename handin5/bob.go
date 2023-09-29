package handin5

import "math/big"

type Bob struct {
	y   int      // Bob input
	d   int      // decode value bob
	F   *big.Int // Encrypted Circuit (missing right type)
	e_y *big.Int // Input encoding Bob
}
