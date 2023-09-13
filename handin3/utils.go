package handin3

// XOR function that returns true if x and y are different, and false if they are the same
func XOR(x bool, y bool) bool {
	return (x || y) && !(x && y)
}

type bloodtype uint8

const (
	ABplus  bloodtype = 0
	ABminus bloodtype = 1
	Bplus   bloodtype = 2
	Bminus  bloodtype = 3
	Aplus   bloodtype = 4
	Aminus  bloodtype = 5
	Oplus   bloodtype = 6
	Ominus  bloodtype = 7
)

func ComputeBloodTypeCompatibility(recipient bloodtype, donor bloodtype) int {
	// Extract the bits from Alice (recipient)
	x := int(recipient)
	y := int(donor)

	x1 := (x >> 2) & 1 // extract 3rd rightmost bit
	x2 := (x >> 1) & 1 // extract 2nd rightmost bit
	x3 := x & 1        // extract rightmost bit

	// Extract the bits from Bob (donor)
	y1 := (y >> 2) & 1 // extract 3rd rightmost bit
	y2 := (y >> 1) & 1 // extract 2nd rightmost bit
	y3 := y & 1        // extract rightmost bit

	return (1 ^ ((1 ^ x1) & y1)) & (1 ^ ((1 ^ x2) & y2)) & (1 ^ ((1 ^ x3) & y3))
}
