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
	m := int(math.Pow(2, float64(8))) // 2^8 = 256. Otherwise it takes too long time to run the test.

	// create cryptoalgorithm, messages and selection bits for algorithms.
	elGamal := elgamal.ElGamal{}
	elGamal.Init()
	selectionBits := utils.RandomSelectionBits(m)
	var messages []*utils.MessagePair
	for i := 0; i < m; i++ {
		msg := utils.MessagePair{
			Message0: utils.RandomBits(l),
			Message1: utils.RandomBits(l),
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

	for iters := 8; iters < 12; iters++ { // k is 128 as we require k <= m
		m := int(math.Pow(2, float64(iters)))

		// create cryptoalgorithm, messages and selection bits for algorithms.
		elGamal := elgamal.ElGamal{}
		elGamal.Init()
		selectionBits := utils.RandomSelectionBits(m)
		var messages []*utils.MessagePair
		for i := 0; i < m; i++ {
			msg := utils.MessagePair{
				Message0: utils.RandomBits(l),
				Message1: utils.RandomBits(l),
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

	for iters := 8; iters < 12; iters++ { // k is 128 as we require k <= m
		m := int(math.Pow(2, float64(iters)))

		// create cryptoalgorithm, messages and selection bits for algorithms.
		elGamal := elgamal.ElGamal{}
		elGamal.Init()
		selectionBits := utils.RandomSelectionBits(m)
		var messages []*utils.MessagePair
		for i := 0; i < m; i++ {
			msg := utils.MessagePair{
				Message0: utils.RandomBits(l),
				Message1: utils.RandomBits(l),
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

	for iters := 8; iters < 12; iters++ { // k is 128 as we require k <= m
		m := int(math.Pow(2, float64(iters)))

		// create cryptoalgorithm, messages and selection bits for algorithms.
		elGamal := elgamal.ElGamal{}
		elGamal.Init()
		selectionBits := utils.RandomSelectionBits(m)
		var messages []*utils.MessagePair
		for i := 0; i < m; i++ {
			msg := utils.MessagePair{
				Message0: utils.RandomBits(l),
				Message1: utils.RandomBits(l),
			}
			messages = append(messages, &msg)
		}

		plaintext := OTExt.OTExtensionProtocolEklundh(k, l, m, selectionBits, messages, elGamal, false)

		// Check if the plaintext is correct
		for i := 0; i < m; i++ {
			if selectionBits[i] == 0 {
				if !bytes.Equal(plaintext[i], messages[i].Message0) {
					t.Errorf("Plaintext is not for selection bit case 0 with plaintext %d and message pair %d, %d for 2^ %d", plaintext[i], messages[i].Message0, messages[i].Message1, iters)
				}
			} else {
				if !bytes.Equal(plaintext[i], messages[i].Message1) {
					t.Errorf("Plaintext is not for selection bit case 1 with plaintext %d and message pair %d, %d for 2^ %d", plaintext[i], messages[i].Message0, messages[i].Message1, iters)
				}
			}
		}
	}

}

func TestOTExtensionProtocolEklundhMultithreaded(t *testing.T) {
	k := 128
	l := 1

	for iters := 8; iters < 12; iters++ { // k is 128 as we require k <= m
		m := int(math.Pow(2, float64(iters)))

		// create cryptoalgorithm, messages and selection bits for algorithms.
		elGamal := elgamal.ElGamal{}
		elGamal.Init()
		selectionBits := utils.RandomSelectionBits(m)
		var messages []*utils.MessagePair
		for i := 0; i < m; i++ {
			msg := utils.MessagePair{
				Message0: utils.RandomBits(l),
				Message1: utils.RandomBits(l),
			}
			messages = append(messages, &msg)
		}

		plaintext := OTExt.OTExtensionProtocolEklundh(k, l, m, selectionBits, messages, elGamal, true)

		// Check if the plaintext is correct
		for i := 0; i < m; i++ {
			if selectionBits[i] == 0 {
				if !bytes.Equal(plaintext[i], messages[i].Message0) {
					t.Errorf("Plaintext is not for selection bit case 0 with plaintext %d and message pair %d, %d for 2^ %d", plaintext[i], messages[i].Message0, messages[i].Message1, iters)
				}
			} else {
				if !bytes.Equal(plaintext[i], messages[i].Message1) {
					t.Errorf("Plaintext is not for selection bit case 1 with plaintext %d and message pair %d, %d for 2^ %d", plaintext[i], messages[i].Message0, messages[i].Message1, iters)
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

			// Additional consistency check. Generate again with the same seed and compare outputs
			output2, _ := utils.PseudoRandomGenerator(seed, length)
			if !equal(output, output2) {
				t.Errorf("Inconsistent outputs for seed %d and length %d", seedVal, length)
			}
		}
	}
}

// equal checks if two slices of byte are equal
func equal(a, b []byte) bool {
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
func TestEklundhTransposeSymmetrical(t *testing.T) {
	for size := 1; size <= 4048; size *= 2 { // Test for different matrix sizes from 1x1 to 16384x16384
		matrix := generateSymmetricMatrix(size)

		expected := utils.TransposeMatrix(matrix)
		result := utils.EklundhTranspose(matrix, false)

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("EklundhTransposeMatrixIterative failed for size %d: got %v, want %v", size, result, expected)
		}
	}
}

func TestEklundhTransposeNonSymmetricalMoreColsThanRows(t *testing.T) {

	rows := 128                                  // rows same as k parameter in OTExtension (most common case)
	for cols := 128; cols <= 262144; cols *= 2 { // Test for different matrix sizes from 128x128 to 128x262144 (2^18). We require rows <= cols (m <= k in OTExtension)
		matrix := generateNonSymmetricMatrix(rows, cols)

		expected := utils.TransposeMatrix(matrix)
		result := utils.EklundhTranspose(matrix, false)

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("EklundhTranspose failed for size %dx%d: got %v, want %v", rows, cols, result, expected)
		}

		if len(result) != cols || len(result[0]) != rows {
			t.Errorf("EklundhTranspose failed in legnth for size %dx%d: got %v, want %v", rows, cols, result, expected)
		}

		for i := 0; i < cols; i++ {
			for j := 0; j < rows; j++ {
				if result[i][j] != expected[i][j] {
					t.Errorf("EklundhTranspose failed for size %dx%d: got %v, want %v", rows, cols, result, expected)
				}
			}
		}

	}
}

// generateNonSymmetricMatrix generates a non-symmetric matrix of size rows x cols.
func generateNonSymmetricMatrix(rows, cols int) [][]byte {
	matrix := make([][]byte, rows)
	for i := range matrix {
		matrix[i] = make([]byte, cols)
		for j := range matrix[i] {
			matrix[i][j] = byte(rand.Intn(256)) // Random byte value
		}
	}
	return matrix
}
