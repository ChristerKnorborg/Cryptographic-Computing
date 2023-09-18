package handin3

import "math/rand"

type Bob struct {
	UVW        []UVW
	x1, x2, x3 int // shares of Alice's input bits from Alice's input x
	y1, y2, y3 int // Bob shares of his input bits from his input y
	e1, e2, e3 int // e-values to be received from Alice
	d1, d2, d3 int // d-values to be received from Alice
	z1, z2, z3 int // z-values to be computed after first layer of AND gates with input bits.
}

// Get random values u, v, w from the dealer and store them in Bob's UVW
func (b *Bob) Init(uvw []UVW) {
	b.UVW = uvw
}

func (b *Bob) TakeInput(y1 int, y2 int, y3 int) (int, int, int) {

	// Generate random shares for Alice from input bits
	y1_a := rand.Intn(2)
	y2_a := rand.Intn(2)
	y3_a := rand.Intn(2)

	// Set Bob's shares
	b.y1 = y1 ^ y1_a
	b.y2 = y2 ^ y2_a
	b.y3 = y3 ^ y3_a

	// Return shares to Alice
	return y1_a, y2_a, y3_a
}

// Bob receives Alice's shares of her input x in bits x1, x2 and x3
func (b *Bob) ReceiveInputShares(x1 int, x2 int, x3 int) {
	b.x1 = x1
	b.x2 = x2
	b.x3 = x3
}

func (b *Bob) Stage1() (int, int, int, int, int, int) {
	b.d1 = b.x1 ^ b.UVW[0].U // Bob masks first bit of his x share:  d1 = x1 ⊕ u1
	b.d2 = b.x2 ^ b.UVW[1].U // Bob masks second bit of his x share: d2 = x2 ⊕ u2
	b.d3 = b.x3 ^ b.UVW[2].U // Bob masks third bit of his x share:  d3 = x3 ⊕ u3

	b.e1 = b.y1 ^ b.UVW[0].V // Bob masks first bit of his y share:  e1 = y1 ⊕ v1
	b.e2 = b.y2 ^ b.UVW[1].V // Bob masks second bit of his y share: e2 = y2 ⊕ v2
	b.e3 = b.y3 ^ b.UVW[2].V // Bob masks third bit of his y share:  e3 = y3 ⊕ v3

	return b.d1, b.d2, b.d3, b.e1, b.e2, b.e3

}

func (b *Bob) Stage2(d1 int, d2 int, d3 int, e1 int, e2 int, e3 int) (int, int) {

	// Bob receives masked d from Alice and unmasks them using his own shares of d
	b.d1 = d1 ^ b.d1
	b.d2 = d2 ^ b.d2
	b.d3 = d3 ^ b.d3

	// Bob receives masked e from Alice and unmasks them using his own shares of e
	b.e1 = e1 ^ b.e1
	b.e2 = e2 ^ b.e2
	b.e3 = e3 ^ b.e3

	// The output share is computed: [z] = [w] ⊕ e & [x] ⊕ d & [y] ⊕ e & d
	// The last step is NOT done for BOB as it is XORing with a constant for the recreated e and d
	b.z1 = b.UVW[0].W ^ (b.e1 & b.x1) ^ (b.d1 & b.y1)
	b.z2 = b.UVW[1].W ^ (b.e2 & b.x2) ^ (b.d2 & b.y2)
	b.z3 = b.UVW[2].W ^ (b.e3 & b.x3) ^ (b.d3 & b.y3)

	// Bob prepares the next AND between z1 and z2
	b.d1 = b.z1 ^ b.UVW[3].U
	b.e1 = b.z2 ^ b.UVW[3].V

	return b.d1, b.e1

}

func (b *Bob) Stage3(d_a int, e_a int) (int, int) {

	// Bob receives masked e and d from Alice and unmasks them using his own shares of e
	b.d1 = b.d1 ^ d_a
	b.e1 = b.e1 ^ e_a

	// The output share is computed: [z] = [w] ⊕ e & [x] ⊕ d & [y] ⊕ e & d
	// The last step is NOT done for BOB as it is XORing with a constant for the recreated e and d
	b.z1 = b.UVW[3].W ^ (b.e1 & b.z1) ^ (b.d1 & b.z2)

	// Bob prepares the next AND between the result of the AND above (saved in z1) and z3 to be used in the next stage.
	b.d1 = b.z1 ^ b.UVW[4].U
	b.e1 = b.z3 ^ b.UVW[4].V

	return b.d1, b.e1
}

func (b *Bob) Stage4(d_a int, e_a int) int {

	// Bob receives masked e from Alice and unmasks them using his own shares of e
	b.d1 = b.d1 ^ d_a
	b.e1 = b.e1 ^ e_a

	// The output share is computed: [z] = [w] ⊕ e & [x] ⊕ d & [y] ⊕ e & d.
	// Notice, the last AND between z1 and z2 is stored in z1.
	// Also notice, the last AND with the recreated e and d values ONLY appear for Alice (e.g. addition with constant)
	b.z1 = b.UVW[4].W ^ (b.e1 & b.z1) ^ (b.d1 & b.z3)

	return b.z1

}
