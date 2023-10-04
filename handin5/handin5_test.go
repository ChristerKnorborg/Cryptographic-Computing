package handin5

import (
	"testing"
)

func TestABplusCanRecieveFromAll(t *testing.T) {
	for donor := Ominus; donor <= ABplus; donor++ {
		if !GarbledCircuit(ABplus, donor) {
			t.Errorf("AB+ should be able to receive from all blood types, failed for donor type: %s", GetBloodTypeName(donor))
		}
	}
}

func TestABplusCanRecieveFromLookUp(t *testing.T) {
	for donor := Ominus; donor <= ABplus; donor++ {
		if !LookUpBloodType(ABplus, donor) {
			t.Errorf("AB+ should be able to receive from all blood types, failed for donor type: %s", GetBloodTypeName(donor))
		}
	}
}

func TestABplusCanRecieveFromAllLookUpComparison(t *testing.T) {
	for donor := Ominus; donor <= ABplus; donor++ {
		if !LookUpBloodType(ABplus, donor) || !GarbledCircuit(ABplus, donor) {
			t.Errorf("AB+ should be able to receive from all blood types, failed for donor type: " + GetBloodTypeName(donor))
		}
	}
}

func TestABminusCanRecieveFromAllMinus(t *testing.T) {
	validDonors := []Bloodtype{ABminus, Aminus, Bminus, Ominus} // AB- can receive from AB-, A-, B-, O-
	for _, donor := range validDonors {
		if !LookUpBloodType(ABminus, donor) || !GarbledCircuit(ABminus, donor) {
			t.Errorf("AB- should be able to receive from AB-, A-, B-, O-, failed for donor type: " + GetBloodTypeName(donor))
		}
	}
}

func TestBplusCanRecieveFromBplusAndOplus(t *testing.T) {
	validDonors := []Bloodtype{Bplus, Bminus, Oplus, Ominus} // B+ can receive from B+, B-, O+, O-
	for _, donor := range validDonors {
		if !LookUpBloodType(Bplus, donor) || !GarbledCircuit(Bplus, donor) {
			t.Errorf("B+ should be able to receive from B+, B-, O+, O-, failed for donor type: " + GetBloodTypeName(donor))
		}
	}
}

func TestBminusCanRecieveFromBminusAndOminus(t *testing.T) {
	validDonors := []Bloodtype{Bminus, Ominus} // B- can receive from B-, O-
	for _, donor := range validDonors {
		if !LookUpBloodType(Bminus, donor) || !GarbledCircuit(Bminus, donor) {
			t.Errorf("B- should be able to receive from B-, O-, failed for donor type: " + GetBloodTypeName(donor))
		}
	}
}

func TestAplusCanRecieveFromAplusAndOplus(t *testing.T) {
	validDonors := []Bloodtype{Aplus, Aminus, Oplus, Ominus} // A+ can receive from A+, A-, O+, O-
	for _, donor := range validDonors {
		if !LookUpBloodType(Aplus, donor) || !GarbledCircuit(Aplus, donor) {
			t.Errorf("A+ should be able to receive from A+, A-, O+, O-, failed for donor type: " + GetBloodTypeName(donor))
		}
	}
}

func TestAminusCanRecieveFromAminusAndOminus(t *testing.T) {
	validDonors := []Bloodtype{Aminus, Ominus} // A- can receive from A-, O-
	for _, donor := range validDonors {
		if !LookUpBloodType(Aminus, donor) || !GarbledCircuit(Aminus, donor) {
			t.Errorf("A- should be able to receive from A-, O-, failed for donor type: " + GetBloodTypeName(donor))
		}
	}
}

func TestOplusCanRecieveFromOplusAndOminus(t *testing.T) {
	validDonors := []Bloodtype{Oplus, Ominus} // O+ can receive from O+, O-
	for _, donor := range validDonors {
		if !LookUpBloodType(Oplus, donor) || !GarbledCircuit(Oplus, donor) {
			t.Errorf("O+ should be able to receive from O+, O-, failed for donor type: " + GetBloodTypeName(donor))
		}
	}
}

func TestOMinusCanRecieveFromOMinus(t *testing.T) {
	if !LookUpBloodType(Ominus, Ominus) || !GarbledCircuit(Ominus, Ominus) {
		t.Errorf("O- should be able to receive from O-, but test failed")
	}
}

// Adding more diverse test case
func TestABplusCannotRecieveFromOMinus(t *testing.T) {
	if !LookUpBloodType(ABplus, Ominus) || !GarbledCircuit(ABplus, Ominus) {
		t.Errorf("AB+ should be able to receive from O-, but test failed")
	}
}

// Adding more diverse test case
func TestOplusCannotRecieveFromABplus(t *testing.T) {
	if LookUpBloodType(Oplus, ABplus) || GarbledCircuit(Oplus, ABplus) {
		t.Errorf("O+ should not be able to receive from AB+, but test passed")
	}
}

// Test that the protocol returns the same result as the lookup table for all possible combinations of blood types
func TestSameOutputLookupTableandBooleanCircuit(t *testing.T) {
	for recipient := Ominus; recipient <= ABplus; recipient++ {
		for donor := Ominus; donor <= ABplus; donor++ {
			lookupResult := LookUpBloodType(recipient, donor)
			otttReseult := GarbledCircuit(recipient, donor)
			if lookupResult != otttReseult {
				t.Fatalf("Mismatch between lookup table and BeDOZa protocol for recipient : " + GetBloodTypeName(recipient) + " and doner: " + GetBloodTypeName(donor))
			}
		}
	}
}
