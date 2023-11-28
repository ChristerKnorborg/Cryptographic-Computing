package project

import (
	"bytes"
	OTBasic "cryptographic-computing/project/OTBasic"
	OTExt "cryptographic-computing/project/OTExtension"
	"cryptographic-computing/project/elgamal"
	utils "cryptographic-computing/project/utils"
	"math"
	"math/big"
	"math/rand"
	"reflect"
	"testing"
)

func TestOTBasicProtocol(t *testing.T) {
	l := 1
	m := int(math.Pow(2, float64(8)))

	// create cryptoalgorithm, messages and selection bits for algorithms.
	elGamal := elgamal.ElGamal{}
	elGamal.Init()
	selectionBits := utils.RandomSelectionBits(m)
	var messages []*utils.MessagePair
	for i := 0; i < m; i++ {
		msg := utils.MessagePair{
			Message0: utils.RandomBytes(l),
			Message1: utils.RandomBytes(l),
		}
		messages = append(messages, &msg)
	}

	plaintext := OTBasic.OTBasicProtocol(l, m, selectionBits, messages, elGamal)

	// Check if the plaintext is correct
	for i := 0; i < m; i++ {
		if selectionBits[i] == 0 {
			if !bytes.Equal(plaintext[i], messages[i].Message0) {
				t.Errorf("Plaintext is not correct")
			}
		} else {
			if !bytes.Equal(plaintext[i], messages[i].Message1) {
				t.Errorf("Plaintext is not correct")
			}
		}
	}

}

func TestOTExtensionProtocol(t *testing.T) {
	k := 128
	l := 1

	for iters := 1; iters < 12; iters++ {
		m := int(math.Pow(2, float64(iters)))

		// create cryptoalgorithm, messages and selection bits for algorithms.
		elGamal := elgamal.ElGamal{}
		elGamal.Init()
		selectionBits := utils.RandomSelectionBits(m)
		var messages []*utils.MessagePair
		for i := 0; i < m; i++ {
			msg := utils.MessagePair{
				Message0: utils.RandomBytes(l),
				Message1: utils.RandomBytes(l),
			}
			messages = append(messages, &msg)
		}

		plaintext := OTExt.OTExtensionProtocol(k, l, m, selectionBits, messages, elGamal)

		// Check if the plaintext is correct
		for i := 0; i < m; i++ {
			if selectionBits[i] == 0 {
				if !bytes.Equal(plaintext[i], messages[i].Message0) {
					t.Errorf("Plaintext is not correct")
				}
			} else {
				if !bytes.Equal(plaintext[i], messages[i].Message1) {
					t.Errorf("Plaintext is not correct")
				}
			}
		}
	}

}

func TestOTExtensionProtocolTranspose(t *testing.T) {
	k := 128
	l := 1

	for iters := 1; iters < 12; iters++ {
		m := int(math.Pow(2, float64(iters)))

		// create cryptoalgorithm, messages and selection bits for algorithms.
		elGamal := elgamal.ElGamal{}
		elGamal.Init()
		selectionBits := utils.RandomSelectionBits(m)
		var messages []*utils.MessagePair
		for i := 0; i < m; i++ {
			msg := utils.MessagePair{
				Message0: utils.RandomBytes(l),
				Message1: utils.RandomBytes(l),
			}
			messages = append(messages, &msg)
		}

		plaintext := OTExt.OTExtensionProtocolTranspose(k, l, m, selectionBits, messages, elGamal)

		// Check if the plaintext is correct
		for i := 0; i < m; i++ {
			if selectionBits[i] == 0 {
				if !bytes.Equal(plaintext[i], messages[i].Message0) {
					t.Errorf("Plaintext is not correct")
				}
			} else {
				if !bytes.Equal(plaintext[i], messages[i].Message1) {
					t.Errorf("Plaintext is not correct")
				}
			}
		}
	}

}

func TestOTExtensionProtocolEklundh(t *testing.T) {
	k := 128
	l := 1

	for iters := 1; iters < 12; iters++ {
		m := int(math.Pow(2, float64(iters)))

		// create cryptoalgorithm, messages and selection bits for algorithms.
		elGamal := elgamal.ElGamal{}
		elGamal.Init()
		selectionBits := utils.RandomSelectionBits(m)
		var messages []*utils.MessagePair
		for i := 0; i < m; i++ {
			msg := utils.MessagePair{
				Message0: utils.RandomBytes(l),
				Message1: utils.RandomBytes(l),
			}
			messages = append(messages, &msg)
		}

		plaintext := OTExt.OTExtensionProtocolEklundh(k, l, m, selectionBits, messages, elGamal)

		// Check if the plaintext is correct
		for i := 0; i < m; i++ {
			if selectionBits[i] == 0 {
				if !bytes.Equal(plaintext[i], messages[i].Message0) {
					t.Errorf("Plaintext is not correct")
				}
			} else {
				if !bytes.Equal(plaintext[i], messages[i].Message1) {
					t.Errorf("Plaintext is not correct")
				}
			}
		}
	}

}

// TestPseudoRandomGeneratorVariety tests the PseudoRandomGenerator function with a variety of seed sizes and output lengths.
func TestPseudoRandomGeneratorVariety(t *testing.T) {
	// Define a range of different seed values and output lengths
	seedValues := []int64{0, 1, 123, 4567, 89012, 345678}
	outputLengths := []int{10, 50, 100, 200, 500}

	for _, seedVal := range seedValues {
		for _, length := range outputLengths {

			seed := big.NewInt(seedVal)
			output, err := utils.PseudoRandomGenerator(seed, length)
			if err != nil {
				t.Errorf("Error returned for seed %d and length %d: %v", seedVal, length, err)
			}
			if len(output) != length {
				t.Errorf("Output length is incorrect for seed %d and length %d. Expected %d, got %d", seedVal, length, length, len(output))
			}

			// Additional consistency check: Generate again with the same seed and compare outputs
			output2, _ := utils.PseudoRandomGenerator(seed, length)
			if !equal(output, output2) {
				t.Errorf("Inconsistent outputs for seed %d and length %d", seedVal, length)
			}
		}
	}
}

// equal checks if two slices of uint8 are equal
func equal(a, b []uint8) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

// generateSymmetricMatrix generates a symmetric matrix of size n x n.
func generateSymmetricMatrix(n int) [][]byte {
	matrix := make([][]byte, n)
	for i := range matrix {
		matrix[i] = make([]byte, n)
		for j := range matrix[i] {
			matrix[i][j] = byte(rand.Intn(256)) // Random byte value
		}
	}
	return matrix
}

// TestEklundhTranspose tests the EklundhTransposeMatrixIterative function.
func TestEklundhTranspose(t *testing.T) {
	for size := 1; size <= 16; size *= 2 {
		matrix := generateSymmetricMatrix(size)

		expected := utils.TransposeMatrix(matrix)
		result := utils.EklundhTransposeMatrix(matrix)

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("EklundhTransposeMatrixIterative failed for size %d: got %v, want %v", size, result, expected)
		}
	}
}
