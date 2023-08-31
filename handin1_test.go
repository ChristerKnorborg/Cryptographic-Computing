package CryptographicComputing

import (
	"testing"
)

// Test for the same return value between the lookup and boolean circuit function
func TestSameOutputLookupTableandBooleanCircuit(t *testing.T) {
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
