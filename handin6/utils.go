package handin6

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

// Bloodtype compatibility lookup table
var Bloodtype_compatibility [8][8]int = [8][8]int{
	{1, 0, 0, 0, 0, 0, 0, 0}, // O- can receive from O-
	{1, 1, 0, 0, 0, 0, 0, 0}, // O+ can receive from O+, O-
	{1, 0, 1, 0, 0, 0, 0, 0}, // B- can receive from B-, O-
	{1, 1, 1, 1, 0, 0, 0, 0}, // B+ can receive from B+, B-, O+, O-
	{1, 0, 0, 0, 1, 0, 0, 0}, // A- can receive from A-, O-
	{1, 1, 0, 0, 1, 1, 0, 0}, // A+ can receive from A+, A-, O+, O-
	{1, 0, 1, 0, 1, 0, 1, 0}, // AB- can receive from AB-, A-, B-, O-
	{1, 1, 1, 1, 1, 1, 1, 1}, // AB+ can receive from everyone
}

// LookUpBloodtype checks if recipient blood type can receive donor blood type using lookup table
func LookUpBloodType(recipient Bloodtype, donor Bloodtype) bool {
	return Bloodtype_compatibility[recipient][donor] == 1
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

// This function extract the three bits from an bloodtype input and returns them as an array
func ExtractBits(n int) [3]int {

	n1 := (n >> 2) & 1 // extract 3rd rightmost bit
	n2 := (n >> 1) & 1 // extract 2nd rightmost bit
	n3 := n & 1        // extract rightmost bit

	return [3]int{n1, n2, n3}
}
