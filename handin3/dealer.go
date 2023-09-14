package handin3

import "math/rand"

type UVW struct {
	U, V, W int
}

type Dealer struct {
	AliceUVW []UVW // alice's u, v, w value pairs for each andGate
	BobUVW   []UVW // bob's u, v, w value pairs for each andGate
}

func (d *Dealer) Init(andGates int) {

	for i := 0; i < andGates; i++ {
		d.GenerateRandomNumbers() // generate random numbers for u and v
	}
}

// The function that generates random numbers u, v and w where w = u * v.
// It makes secret shares for the numbers, one share for Alice and one share for Bob.
func (d *Dealer) GenerateRandomNumbers() {

	// Generate real u, v, w
	u := rand.Intn(2)
	v := rand.Intn(2)
	w := u & v // w = u * v

	// Generate random shares for Alice
	u_A := rand.Intn(2)
	v_A := rand.Intn(2)
	w_A := rand.Intn(2)

	// Calculate the complementing shares for Bob
	u_B := u - u_A
	v_B := v - v_A
	w_B := w - w_A

	// Make sure to mod by 2 to keep it a single bit.
	// The + 2 is to make sure it's positive since golang's modulo operator does not support negative
	u_B = (u_B + 2) % 2
	v_B = (v_B + 2) % 2
	w_B = (w_B + 2) % 2

	// Append u, v, w to alice slice and bob slice
	d.AliceUVW = append(d.AliceUVW, UVW{u_A, v_A, w_A})
	d.BobUVW = append(d.BobUVW, UVW{u_B, v_B, w_B})
}

// The function that returns a single tuple of the numbers u, v and w for Alice
func (d *Dealer) GetAliceUVW() []UVW {
	return d.AliceUVW // get first element of aliceNumbers slice
}

// The function that returns a single tuple of the numbers u, v and w for Bob
func (d *Dealer) GetBobUVW() []UVW {
	return d.BobUVW // get first element of bobNumbers slice
}
