package handin3

import "math/rand"

type Alice struct {
	UVW        []UVW
	x1, x2, x3 int // Alice shares of her input bits from her input x
	y1, y2, y3 int // shares of Bob's input bits from Bob's input y
	e1, e2, e3 int // e-values to be received from Bob
	d1, d2, d3 int // d-values to be received from Bob
	z1, z2, z3 int // z-values to be computed after first layer of AND gates with input bits.
}

func (a *Alice) Init(uvw []UVW) {
	a.UVW = uvw
}

// Generate shares of Alice's input x.
func (a *Alice) TakeInput(x1 int, x2 int, x3 int) (int, int, int) {

	// Generate random shares for Alice from input bits
	x1_b := rand.Intn(2)
	x2_b := rand.Intn(2)
	x3_b := rand.Intn(2)

	// Set Alice's shares
	a.x1 = x1 ^ x1_b
	a.x2 = x2 ^ x2_b
	a.x3 = x3 ^ x3_b

	// Return shares to Bob
	return x1_b, x2_b, x3_b
}

// Alice receives Bob's shares of his input y in bits y1, y2 and y3
func (a *Alice) ReceiveInputShares(y1 int, y2 int, y3 int) {
	a.y1 = y1
	a.y2 = y2
	a.y3 = y3
}

func (a *Alice) Stage1() (int, int, int, int, int, int) {

	a.x1 = a.x1 ^ 1 // Alice negates her first input bit by XORing with constant 1
	a.x2 = a.x2 ^ 1 // Alice negates her second input bit by XORing with constant 1
	a.x3 = a.x3 ^ 1 // Alice negates her third input bit by XORing with constant 1

	a.d1 = a.x1 ^ a.UVW[0].U // Alice masks first bit of her x share:  d1 = x1 ⊕ u1
	a.d2 = a.x2 ^ a.UVW[1].U // Alice masks second bit of her x share: d2 = x2 ⊕ u2
	a.d3 = a.x3 ^ a.UVW[2].U // Alice masks third bit of her x share:  d3 = x3 ⊕ u3

	a.e1 = a.y1 ^ a.UVW[0].V // Alice masks first bit of her y share:  e1 = y1 ⊕ v1
	a.e2 = a.y2 ^ a.UVW[1].V // Alice masks second bit of her y share: e2 = y2 ⊕ v2
	a.e3 = a.y3 ^ a.UVW[2].V // Alice masks third bit of her y share:  e3 = y3 ⊕ v3

	return a.d1, a.d2, a.d3, a.e1, a.e2, a.e3
}

func (a *Alice) Stage2(d1 int, d2 int, d3 int, e1 int, e2 int, e3 int) (int, int) {

	// Alice receives masked d from Bob and unmasks them using her own shares of d
	a.d1 = d1 ^ a.d1
	a.d2 = d2 ^ a.d2
	a.d3 = d3 ^ a.d3

	// Alice receives masked e from Bob and unmasks them using her own shares of e
	a.e1 = e1 ^ a.e1
	a.e2 = e2 ^ a.e2
	a.e3 = e3 ^ a.e3

	// The output share is computed: [z] = [w] ⊕ e & [x] ⊕ d & [y] ⊕ e & d
	// Notice, the last AND with the recreated e and d values ONLY appear for Alice (e.g. addition with constant)
	a.z1 = a.UVW[0].W ^ (a.e1 & a.x1) ^ (a.d1 & a.y1) ^ (a.e1 & a.d1)
	a.z2 = a.UVW[1].W ^ (a.e2 & a.x2) ^ (a.d2 & a.y2) ^ (a.e2 & a.d2)
	a.z3 = a.UVW[2].W ^ (a.e3 & a.x3) ^ (a.d3 & a.y3) ^ (a.e3 & a.d3)

	// Alice negates the output of the AND gates - Only Alice does this as constant 1 is only added to Alice's shares
	a.z1 = a.z1 ^ 1
	a.z2 = a.z2 ^ 1
	a.z3 = a.z3 ^ 1

	// Alice prepares the next AND between z1 and z2 (z3 is not used until the next layer)
	a.d1 = a.z1 ^ a.UVW[3].U
	a.e1 = a.z2 ^ a.UVW[3].V

	return a.d1, a.e1
}

func (a *Alice) Stage3(d_b int, e_b int) (int, int) {

	// Alice receives masked e and d from Bob and unmasks them using her own shares of e
	a.d1 = a.d1 ^ d_b
	a.e1 = a.e1 ^ e_b

	// The output share is computed: [z] = [w] ⊕ e & [x] ⊕ d & [y] ⊕ e & d
	// Notice, the last AND with the recreated e and d values ONLY appear for Alice (e.g. addition with constant)
	a.z1 = a.UVW[3].W ^ (a.e1 & a.z1) ^ (a.d1 & a.z2) ^ (a.e1 & a.d1)

	// Alice prepares the next AND between the result of the AND above (saved in z1) and z3 to be used in the next stage.
	a.d1 = a.z1 ^ a.UVW[4].U
	a.e1 = a.z3 ^ a.UVW[4].V

	return a.d1, a.e1
}

func (a *Alice) Stage4(d_b int, e_b int) int {
	// Alice receives masked e and d from Bob and unmasks them using her own shares of e
	a.d1 = a.d1 ^ d_b
	a.e1 = a.e1 ^ e_b

	// The output share is computed: [z] = [w] ⊕ e & [x] ⊕ d & [y] ⊕ e & d.
	// Notice, the result from the AND between z1 and z2 is saved in z1.
	// Also notice, the last AND with the recreated e and d values ONLY appear for Alice (e.g. addition with constant)
	a.z1 = a.UVW[4].W ^ (a.e1 & a.z1) ^ (a.d1 & a.z3) ^ (a.e1 & a.d1)

	return a.z1
}
