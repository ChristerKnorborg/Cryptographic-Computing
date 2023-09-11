package handin3

import "math/rand"


type UVW struct {
	u, v, w int
}

type Dealer struct {
	AliceUVW UVW[] // alice's u, v, w value pairs for each andGate
	AliceUVW UVW[] // bob's u, v, w value pairs for each andGate
}

func (d *Dealer) Init(int andGates) {

	var andGates [](int, int) // create empty slice of andGates

	for i := 0; i < andGates; i++ {
		GenerateRandomNumbers() // generate random numbers for u and v
	}
}

// The function that generates random numbers u, v and w where w = u * v. 
// Makes two sets of numbers, one for Alice and one for Bob.
func (d *Dealer) GenerateRandomNumbers(){
	
	u_A := rand.Intn(2) // generate random bit for u for Alice
	v_A := rand.Intn(2) // generate random bit for v for Alice
	w_A := u_A && v_A // generate w for Alice where w = u * v

	u_B := rand.Intn(2) // generate random bit for u for Bob
	v_B := rand.Intn(2) // generate random bit for v for Bob
	w_B := u_B && v_B // generate w for Bob where w = u * v

	d.aliceNumbers = append(d.alice, (u_A, v_A, w_A)) // append u, v, w to alice slice
	d.bobNumbers = append(d.bob, (u_B, v_B, w_B)) // append u, v, w to bob slice
}

// The function that returns a single tuple of the numbers u, v and w for Alice
func (d *Dealer) getAliceUVW() [](int, int, int) {
	return d.AliceUVW // get first element of aliceNumbers slice
}

// The function that returns a single tuple of the numbers u, v and w for Bob
func (d *Dealer) getBobUVW() [](int, int, int) {
	return d.BobUVW // get first element of bobNumbers slice
}



