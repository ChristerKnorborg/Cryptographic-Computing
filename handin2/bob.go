package handin2

type Bob struct {
	y        int
	matrix_b [8][8]bool
	s        int
	u        int
}

func (b *Bob) Init(y int, matrix [8][8]bool, s int) {
	b.y = y             // Bob inputs his own y
	b.matrix_b = matrix // Bob get his matrix from the dealer
	b.s = s             // Bob get his s from the dealer
}

func (b *Bob) Send() (int, bool) {
	v := (int(b.y) + b.s) % 8 // Bob computes v = y + s mod n,
	z_B := b.matrix_b[b.u][v] // and z_B = M_B[u, v] and sends (v, z_B) to Alice.
	return v, z_B
}

// Method that immitates Bob receiving data from Alice on a channel between Alice and Bob
func (b *Bob) Receive(u int) {
	b.u = u // Bob gets u = x + r mod n from Alice
}
