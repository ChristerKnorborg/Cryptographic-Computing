package CryptographicComputing

import (
	"testing"
)

// TestHelloName calls greetings.Hello with a name, checking
// for a valid return value.
func TestHelloName(t *testing.T) {
	for recipient := ABplus; recipient <= Ominus; recipient++ {
		for donor := ABplus; donor <= Ominus; donor++ {
			lookupResult := LookUpBloodType(recipient, donor)
			formulaResult := BooleanFormula(recipient, donor)
			if lookupResult != formulaResult {
				t.Fatalf("Mismatch for lookupresult = %q, and formularesult = %p", recipient, &donor)
			}
		}
	}
}
