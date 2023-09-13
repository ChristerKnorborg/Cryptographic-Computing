package handin3

import "math/rand"

type Alice struct {
	UVW        []UVW
	x1, x2, x3 int // Alice shares of her input bits from her input x
	y1, y2, y3 int // shares of Bob's input bits from Bob's input y
	e1, e2, e3 int // e-values to be received from Bob
	d1, d2, d3 int // d-values to be received from Bob
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
func (a *Alice) ReceiveInput(y1 int, y2 int, y3 int) {
	a.y1 = y1
	a.y2 = y2
	a.y3 = y3
}

func (a *Alice) SendValues() (int, int, int, int, int, int) {
	return a.d1, a.d2, a.d3, a.e1, a.e2, a.e3
}

func (a *Alice) ReceiveValues(d1 int, d2 int, d3 int, e1 int, e2 int, e3 int) {
	a.d1 = d1
	a.d2 = d2
	a.d3 = d3
	a.e1 = e1
	a.e2 = e2
	a.e3 = e3
}
