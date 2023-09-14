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
	x1_a := rand.Intn(2)
	x2_a := rand.Intn(2)
	x3_a := rand.Intn(2)

	// Set Alice's shares
	a.x1 = x1_a
	a.x2 = x2_a
	a.x3 = x3_a

	// Calculate the complementing shares for Bob
	x1_b := x1 - x1_a
	x2_b := x2 - x2_a
	x3_b := x3 - x3_a

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

	a.x1 = a.x1 ^ 1
	a.x2 = a.x2 ^ 1
	a.x3 = a.x3 ^ 1

	d1_a := a.x1 ^ a.UVW[0].U // Alice masks first bit of her x share:  d1 = x1 ⊕ u1
	d2_a := a.x2 ^ a.UVW[1].U // Alice masks second bit of her x share: d2 = x2 ⊕ u2
	d3_a := a.x3 ^ a.UVW[2].U // Alice masks third bit of her x share:  d3 = x3 ⊕ u3

	e1_a := a.y1 ^ a.UVW[0].V // Alice masks first bit of her y share:  e1 = y1 ⊕ v1
	e2_a := a.y2 ^ a.UVW[1].V // Alice masks second bit of her y share: e2 = y2 ⊕ v2
	e3_a := a.y3 ^ a.UVW[2].V // Alice masks third bit of her y share:  e3 = y3 ⊕ v3

	return d1_a, d2_a, d3_a, e1_a, e2_a, e3_a
}

func (a *Alice) Stage2(d1 int, d2 int, d3 int, e1 int, e2 int, e3 int) {

	// Alice receives masked d from Bob and unmasks them using her own shares of d
	a.d1 = d1 ^ a.d1
	a.d2 = d2 ^ a.d2
	a.d3 = d3 ^ a.d3

	// Alice receives masked e from Bob and unmasks them using her own shares of e
	a.e1 = e1 ^ a.e1
	a.e2 = e2 ^ a.e2
	a.e3 = e3 ^ a.e3

	// The output share is computed: [z] = [w] ⊕ e & [x] ⊕ d & [y] ⊕ e & d
	a.z1 = a.UVW[0].W ^ a.e1&a.x1 ^ a.d1&a.y1 ^ a.e1&a.d1
	a.z2 = a.UVW[1].W ^ a.e2&a.x2 ^ a.d2&a.y2 ^ a.e2&a.d2
	a.z3 = a.UVW[2].W ^ a.e3&a.x3 ^ a.d3&a.y3 ^ a.e3&a.d3

	// Alice negates the output of the AND gates - Only Alice does this as constant 1 is only added to Alice's shares
	a.z1 = a.z1 ^ 1
	a.z2 = a.z2 ^ 1
	a.z3 = a.z3 ^ 1

}

func (a *Alice) MaskZ1AndZ2() (int, int) {
	// Alice masks the output of the AND gates
	e := a.z1 ^ a.UVW[3].U
	d := a.z2 ^ a.UVW[3].V

	return e, d
}

func (a *Alice) ReceiveMasked(e int, d int) {
	a.e1 = e ^ a.e1
	a.d1 = d ^ a.d1

}
