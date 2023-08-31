package CryptographicComputing

import (
	"testing"
)

func TestABplusCanRecieveFromAll(t *testing.T) {
	for donor := ABplus; donor <= Ominus; donor++ {
		if !LookUpBloodType(ABplus, donor) || !BooleanFormula(ABplus, donor) {
			t.Errorf("AB+ should be able to receive from all blood types, failed for donor type: %d", donor)
		}
	}
}

func TestABminusCanRecieveFromABminusAndOminus(t *testing.T) {
	validDonors := []bloodtype{ABminus, Aminus, Bminus, Ominus} // AB- can receive from AB-, A-, B-, O-
	for _, donor := range validDonors {
		if !LookUpBloodType(ABminus, donor) || !BooleanFormula(ABminus, donor) {
			t.Errorf("AB- should be able to receive from AB-, A-, B-, O-, failed for donor type: %d", donor)
		}
	}
}

func TestBplusCanRecieveFromBplusAndOplus(t *testing.T) {
	validDonors := []bloodtype{Bplus, Bminus, Oplus, Ominus} // B+ can receive from B+, B-, O+, O-
	for _, donor := range validDonors {
		if !LookUpBloodType(Bplus, donor) || !BooleanFormula(Bplus, donor) {
			t.Errorf("B+ should be able to receive from B+, B-, O+, O-, failed for donor type: %d", donor)
		}
	}
}

func TestBminusCanRecieveFromBminusAndOminus(t *testing.T) {
	validDonors := []bloodtype{Bminus, Ominus} // B- can receive from B-, O-
	for _, donor := range validDonors {
		if !LookUpBloodType(Bminus, donor) || !BooleanFormula(Bminus, donor) {
			t.Errorf("B- should be able to receive from B-, O-, failed for donor type: %d", donor)
		}
	}
}

func TestAplusCanRecieveFromAplusAndOplus(t *testing.T) {
	validDonors := []bloodtype{Aplus, Aminus, Oplus, Ominus} // A+ can receive from A+, A-, O+, O-
	for _, donor := range validDonors {
		if !LookUpBloodType(Aplus, donor) || !BooleanFormula(Aplus, donor) {
			t.Errorf("A+ should be able to receive from A+, A-, O+, O-, failed for donor type: %d", donor)
		}
	}
}

func TestAminusCanRecieveFromAminusAndOminus(t *testing.T) {
	validDonors := []bloodtype{Aminus, Ominus} // A- can receive from A-, O-
	for _, donor := range validDonors {
		if !LookUpBloodType(Aminus, donor) || !BooleanFormula(Aminus, donor) {
			t.Errorf("A- should be able to receive from A-, O-, failed for donor type: %d", donor)
		}
	}
}

func TestOplusCanRecieveFromOplusAndOminus(t *testing.T) {
	validDonors := []bloodtype{Oplus, Ominus} // O+ can receive from O+, O-
	for _, donor := range validDonors {
		if !LookUpBloodType(Oplus, donor) || !BooleanFormula(Oplus, donor) {
			t.Errorf("O+ should be able to receive from O+, O-, failed for donor type: %d", donor)
		}
	}
}

func TestOMinusCanRecieveFromOMinus(t *testing.T) {
	if !LookUpBloodType(Ominus, Ominus) || !BooleanFormula(Ominus, Ominus) {
		t.Errorf("O- should be able to receive from O-, but test failed")
	}
}

// Adding more diverse test case
func TestABplusCannotRecieveFromOMinus(t *testing.T) {
	if !LookUpBloodType(ABplus, Ominus) || !BooleanFormula(ABplus, Ominus) {
		t.Errorf("AB+ should be able to receive from O-, but test failed")
	}
}

// Adding more diverse test case
func TestOplusCannotRecieveFromABplus(t *testing.T) {
	if LookUpBloodType(Oplus, ABplus) || BooleanFormula(Oplus, ABplus) {
		t.Errorf("O+ should not be able to receive from AB+, but test passed")
	}
}

// Test for the same return value between the lookup and boolean circuit function for all possible combinations of blood types
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
