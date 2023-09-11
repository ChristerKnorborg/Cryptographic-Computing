package handin3


func XORWithConstantA(int x_A, int x_B, int c) {
	z_A := XOR(x_A, c)
	z_B := x_B
	z := XOR(z_A, z_B)
	return z

}


func ANDWithConstant(int x_A, int x_B, int c) {
	z_A := x_A && c
	z_B := x_B && c
	z := XOR(z_A, z_B)
	return z
}


func XOROfTwoWires(int x_A, int y_A, int x_B, int y_B) {
	z_A := XOR(x_A, y_A)
	z_B := XOR(x_B, y_B)
	z := XOR(z_A, z_B)
	return z
}

func ANDOfTwoWires(int x_A, int y_A, int u, int x_B, int y_B, int v) {

	w := v && u // [w] = [v] ∧ [u]

	x = XOR(x_A, x_B) // [x] = [x_A] ⊕ [x_B]
	d := XOR(x, u) // [d] = [x] ⊕ [u]

	y := XOR(y_A, y_B) // [y] = [y_A] ⊕ [y_B]
	e := XOR(y, v)  // [e] = [y] ⊕ [v]


	z := XOR(w, e) && XOR(x, d) && XOR(y, e) && d // [z]   =   [w] ⊕ e ∧ [x] ⊕ d ∧ [y] ⊕ e ∧ d   =   uv ⊕ (xy ⊕ vx) ⊕ (xy ⊕ uy) ⊕ (xy ⊕ vx ⊕ uy ⊕ uv)   =   xy
	return z



}

