package handin2

type Alice struct {
	x        int
	matrix_a [8][8]bool
	u        int
	r        int
	v        int
	z_B      bool
}

func (a *Alice) Init(x int, matrix [8][8]bool, r int) {
	a.x = x             // Alice inputs her own x
	a.matrix_a = matrix // Alice get her matrix from the dealer
	a.r = r             // Alice get her r from the dealer

}

// Method that immitates Alice sending data to Bob on a channel between Alice and Bob
func (a *Alice) Send() int {
	u := (int(a.x) + a.r) % 8 // Alice computes u = x + r mod n and sends u to Bob
	a.u = u                   // Alice also stores u (to send to Bob later when immiating sending data on a channel)
	return u
}

// Method that immitates Alice receiving data from Bob on a channel between Alice and Bob
func (a *Alice) Receive(v int, z_B bool) {
	a.v = v     // Alice gets v from Bob
	a.z_B = z_B // Alice gets z_B from Bob
}

func (a *Alice) ComputeOutput() bool {
	z := XOR(a.matrix_a[a.u][a.v], a.z_B) // Alice outputs z = M_A[u, v] ⊕ z_B. Which is equal to f(x,y)
	return z
}
