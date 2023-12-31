package handin3

type Bloodtype uint8

const (
	Ominus  Bloodtype = 0
	Oplus   Bloodtype = 1
	Bminus  Bloodtype = 2
	Bplus   Bloodtype = 3
	Aminus  Bloodtype = 4
	Aplus   Bloodtype = 5
	ABminus Bloodtype = 6
	ABplus  Bloodtype = 7
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
	{true, false, false, false, false, false, false, false}, // O- can receive from O-
	{true, true, false, false, false, false, false, false},  // O+ can receive from O+, O-
	{true, false, true, false, false, false, false, false},  // B- can receive from B-, O-
	{true, true, true, true, false, false, false, false},    // B+ can receive from B+, B-, O+, O-
	{true, false, false, false, true, false, false, false},  // A- can receive from A-, O-
	{true, true, false, false, true, true, false, false},    // A+ can receive from A+, A-, O+, O-
	{true, false, true, false, true, false, true, false},    // AB- can receive from AB-, A-, B-, O-
	{true, true, true, true, true, true, true, true},        // AB+ can receive from everyone
}

// LookUpBloodtype checks if recipient blood type can receive donor blood type using lookup table
func LookUpBloodType(recipient Bloodtype, donor Bloodtype) bool {
	return Bloodtype_compatibility[recipient][donor]
}

func GetBloodTypeName(bType Bloodtype) string {
	switch bType {
	case ABplus:
		return "AB+"
	case ABminus:
		return "AB-"
	case Bplus:
		return "B+"
	case Bminus:
		return "B-"
	case Aplus:
		return "A+"
	case Aminus:
		return "A-"
	case Oplus:
		return "O+"
	case Ominus:
		return "O-"
	default:
		return "Unknown"
	}
}
