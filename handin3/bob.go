package handin3

import "math/rand"

type Bob struct {
	UVW        []UVW
	x1, x2, x3 int // shares of Alice's input bits from Alice's input x
	y1, y2, y3 int // Bob shares of his input bits from his input y
	e1, e2, e3 int // e-values to be received from Alice
	d1, d2, d3 int // d-values to be received from Alice
}

func (b *Bob) Init(uvw []UVW) {
	b.UVW = uvw
}

func (b *Bob) TakeInput(y1 int, y2 int, y3 int) (int, int, int) {

	// Generate random shares for Alice from input bits
	y1_b := rand.Intn(2)
	y2_b := rand.Intn(2)
	y3_b := rand.Intn(2)

	// Set Bob's shares
	b.y1 = y1_b
	b.y2 = y2_b
	b.y3 = y3_b

	// Calculate the complementing shares for Alice
	y1_a := y1 - y1_b
	y2_a := y2 - y2_b
	y3_a := y3 - y3_b

	// Return shares to Alice
	return y1_a, y2_a, y3_a
}

// Bob receives Alice's shares of her input x in bits x1, x2 and x3
func (b *Bob) ReceiveInput(x1 int, x2 int, x3 int) {
	b.x1 = x1
	b.x2 = x2
	b.x3 = x3
}

func (b *Bob) SendValues() (int, int, int, int, int, int) {
	return b.d1, b.d2, b.d3, b.e1, b.e2, b.e3
}

func (b *Bob) ReceiveValues(d1 int, d2 int, d3 int, e1 int, e2 int, e3 int) {
	b.d1 = d1
	b.d2 = d2
	b.d3 = d3
	b.e1 = e1
	b.e2 = e2
	b.e3 = e3
}
