package handin3

import (
	"testing"
)

func TestABplusCanRecieveFromAll(t *testing.T) {
	for donor := ABplus; donor <= Ominus; donor++ {
		if !LookUpBloodType(ABplus, donor) || !ComputeBeDOZaBloodTypeCompatability(ABplus, donor) {
			t.Errorf("AB+ should be able to receive from all blood types, failed for donor type: %d", donor)
		}
	}
}

func TestABminusCanRecieveFromABminusAndOminus(t *testing.T) {
	validDonors := []bloodtype{ABminus, Aminus, Bminus, Ominus} // AB- can receive from AB-, A-, B-, O-
	for _, donor := range validDonors {
		if !LookUpBloodType(ABminus, donor) || !ComputeBeDOZaBloodTypeCompatability(ABminus, donor) {
			t.Errorf("AB- should be able to receive from AB-, A-, B-, O-, failed for donor type: %d", donor)
		}
	}
}

func TestBplusCanRecieveFromBplusAndOplus(t *testing.T) {
	validDonors := []bloodtype{Bplus, Bminus, Oplus, Ominus} // B+ can receive from B+, B-, O+, O-
	for _, donor := range validDonors {
		if !LookUpBloodType(Bplus, donor) || !ComputeBeDOZaBloodTypeCompatability(Bplus, donor) {
			t.Errorf("B+ should be able to receive from B+, B-, O+, O-, failed for donor type: %d", donor)
		}
	}
}

func TestBminusCanRecieveFromBminusAndOminus(t *testing.T) {
	validDonors := []bloodtype{Bminus, Ominus} // B- can receive from B-, O-
	for _, donor := range validDonors {
		if !LookUpBloodType(Bminus, donor) || !ComputeBeDOZaBloodTypeCompatability(Bminus, donor) {
			t.Errorf("B- should be able to receive from B-, O-, failed for donor type: %d", donor)
		}
	}
}

func TestAplusCanRecieveFromAplusAndOplus(t *testing.T) {
	validDonors := []bloodtype{Aplus, Aminus, Oplus, Ominus} // A+ can receive from A+, A-, O+, O-
	for _, donor := range validDonors {
		if !LookUpBloodType(Aplus, donor) || !ComputeBeDOZaBloodTypeCompatability(Aplus, donor) {
			t.Errorf("A+ should be able to receive from A+, A-, O+, O-, failed for donor type: %d", donor)
		}
	}
}

func TestAminusCanRecieveFromAminusAndOminus(t *testing.T) {
	validDonors := []bloodtype{Aminus, Ominus} // A- can receive from A-, O-
	for _, donor := range validDonors {
		if !LookUpBloodType(Aminus, donor) || !ComputeBeDOZaBloodTypeCompatability(Aminus, donor) {
			t.Errorf("A- should be able to receive from A-, O-, failed for donor type: %d", donor)
		}
	}
}

func TestOplusCanRecieveFromOplusAndOminus(t *testing.T) {
	validDonors := []bloodtype{Oplus, Ominus} // O+ can receive from O+, O-
	for _, donor := range validDonors {
		if !LookUpBloodType(Oplus, donor) || !ComputeBeDOZaBloodTypeCompatability(Oplus, donor) {
			t.Errorf("O+ should be able to receive from O+, O-, failed for donor type: %d", donor)
		}
	}
}

func TestOMinusCanRecieveFromOMinus(t *testing.T) {
	if !LookUpBloodType(Ominus, Ominus) || !ComputeBeDOZaBloodTypeCompatability(Ominus, Ominus) {
		t.Errorf("O- should be able to receive from O-, but test failed")
	}
}

// Adding more diverse test case
func TestABplusCannotRecieveFromOMinus(t *testing.T) {
	if !LookUpBloodType(ABplus, Ominus) || !ComputeBeDOZaBloodTypeCompatability(ABplus, Ominus) {
		t.Errorf("AB+ should be able to receive from O-, but test failed")
	}
}

// Adding more diverse test case
func TestOplusCannotRecieveFromABplus(t *testing.T) {
	if LookUpBloodType(Oplus, ABplus) || ComputeBeDOZaBloodTypeCompatability(Oplus, ABplus) {
		t.Errorf("O+ should not be able to receive from AB+, but test passed")
	}
}

// Test that the OTTT protocol returns the same result as the lookup table for all possible combinations of blood types
func TestSameOutputLookupTableandBooleanCircuit(t *testing.T) {
	for recipient := ABplus; recipient <= Ominus; recipient++ {
		for donor := ABplus; donor <= Ominus; donor++ {
			lookupResult := LookUpBloodType(recipient, donor)
			otttReseult := ComputeBeDOZaBloodTypeCompatability(recipient, donor)
			if lookupResult != otttReseult {
				t.Fatalf("Mismatch for lookupresult = %q, and formularesult = %p", recipient, &donor)
			}
		}
	}
}
