package handin3

func XORWithConstant(x int, c int) int {
	return x ^ c

}

func ANDWithConstant(x_A int, x_B int, c int) int {
	z_A := x_A & c
	z_B := x_B & c
	z := z_A ^ z_B
	return z
}

func XOROfTwoWires(x_A int, y_A int, x_B int, y_B int) int {
	z_A := x_A ^ y_A
	z_B := x_B ^ y_B
	z := z_A ^ z_B
	return z
}

func ANDOfTwoWires(x_A int, y_A int, u int, x_B int, y_B int, v int) int {

	w := v & u // [w] = [v] ∧ [u]

	x := x_A ^ x_B // [x] = [x_A] ⊕ [x_B]
	d := x ^ u     // [d] = [x] ⊕ [u]

	y := y_A ^ y_B // [y] = [y_A] ⊕ [y_B]
	e := y ^ v     // [e] = [y] ⊕ [v]

	z := (w ^ e) & (x ^ d) & (y ^ e) & d // [z]   =   [w] ⊕ e ∧ [x] ⊕ d ∧ [y] ⊕ e ∧ d   =   uv ⊕ (xy ⊕ vx) ⊕ (xy ⊕ uy) ⊕ (xy ⊕ vx ⊕ uy ⊕ uv)   =   xy
	return z

}
