package handin3

// XOR function that returns true if x and y are different, and false if they are the same
func XOR(x bool, y bool) bool {
	return (x || y) && !(x && y)
}

type Bloodtype uint8

const (
	ABplus  Bloodtype = 0
	ABminus Bloodtype = 1
	Bplus   Bloodtype = 2
	Bminus  Bloodtype = 3
	Aplus   Bloodtype = 4
	Aminus  Bloodtype = 5
	Oplus   Bloodtype = 6
	Ominus  Bloodtype = 7
)

func ComputeBloodtypeCompatibility(recipient Bloodtype, donor Bloodtype) int {
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

// Bloodtype compatibility lookup table
var Bloodtype_compatibility [8][8]bool = [8][8]bool{
	{true, true, true, true, true, true, true, true},        // AB+
	{false, true, false, true, false, true, false, true},    // AB-
	{false, false, true, true, false, false, true, true},    // B+
	{false, false, false, true, false, false, false, true},  // B-
	{false, false, false, false, true, true, true, true},    // A+
	{false, false, false, false, false, true, false, true},  // A-
	{false, false, false, false, false, false, true, true},  // O+
	{false, false, false, false, false, false, false, true}, // O-
}

// LookUpBloodtype checks if recipient blood type can receive donor blood type using lookup table
func LookUpBloodType(recipient Bloodtype, donor Bloodtype) bool {
	return Bloodtype_compatibility[recipient][donor]
}
