package utils

// Import your ElGamal package
import (
	"crypto/rand"
	"crypto/sha256"
	"cryptographic-computing/project/elgamal"
	"fmt"
	"log"
	"math/big"
	mathRand "math/rand"
	"sync"
	"time"
)

// Struct to store messages M0 and M1
type MessagePair struct {
	Message0 []byte
	Message1 []byte
}

// Struct to store public keys for Oblivious transfer. Each public key pair consists of 2 public keys -
// one real key and one oblivious key.
type PublicKeyPair struct {
	MessageKey0 *big.Int
	MessageKey1 *big.Int
}

// Struct to store ciphertexts of two messages M0 and M1
type CiphertextPair struct {
	Ciphertext0 *elgamal.Ciphertext
	Ciphertext1 *elgamal.Ciphertext
}

// Struct to store seeds for the OTExtension protocol in the initial phase.
type Seed struct {
	Seed0 *big.Int
	Seed1 *big.Int
}

// Struct to store ciphertexts in the OTExtension protocol in the final phase.
type ByteCiphertextPair struct {
	Y0 []byte
	Y1 []byte
}

// PseudoRandomGenerator generates a pseudo-random bit string of a given bit length using sha256 and a seed.
func PseudoRandomGenerator(seed *big.Int, bitLength int) ([]byte, error) {
	if bitLength <= 0 {
		return nil, fmt.Errorf("bitLength must be positive")
	}
	output := make([]byte, 0, bitLength)

	// Convert seed to a byte slice
	seedBytes := seed.Bytes()

	for len(output) < bitLength {
		hash := sha256.Sum256(seedBytes)
		for _, b := range hash[:] {
			for i := 0; i < 8 && len(output) < bitLength; i++ {
				bit := (b >> (7 - i)) & 1 // Extract each bit
				output = append(output, byte(bit))
			}
		}

		// Increment the seed
		seed = new(big.Int).Add(seed, big.NewInt(1))
		seedBytes = seed.Bytes()
	}

	return output, nil
}

// Hash creates a hash of the input data with a specified bit length using sha256.
func Hash(originalData []byte, length int) []byte {

	// Create a copy of the original data to avoid modifying it
	data := make([]byte, len(originalData))
	copy(data, originalData)

	var byteLength int
	if length <= 8 {
		byteLength = 1
	} else {
		byteLength = (length + 7) / 8 // Calculate the number of bytes needed to represent the given number of bits
	}

	fullHash := make([]byte, 0, byteLength)

	// SHA-256 produces a hash of 32 bytes (256 bits)
	// We keep hashing and concatenating until we reach the desired byte length
	for len(fullHash) < byteLength {
		hash := sha256.Sum256(data)
		fullHash = append(fullHash, hash[:]...)

		// We modify data slightly for the next iteration to produce a different hash by appending a byte (the length of the hash in this case)
		data = append(data, byte(len(fullHash)))
	}
	// Truncate the hash to the exact byte length required
	returnHash := fullHash[:byteLength]

	// Set remaining bits to 0 if length is not multiple of 8
	remainingBits := length % 8
	if remainingBits != 0 {
		returnHash[byteLength-1] &= (1 << remainingBits) - 1
	}
	return returnHash
}

// PrintBinaryString prints the binary representation of a []byte slice.
func PrintBinaryString(bytes []byte) {
	binaryString := ""
	for _, b := range bytes {
		binaryString += fmt.Sprintf("%08b", b) // Convert each byte to an 8-bit binary string
	}
	fmt.Println("As binary string:", binaryString)
}

// RandomBits generates a slice of random bytes of a given bit length
func RandomBits(length int) []byte {

	var numBytes int
	if length <= 8 {
		numBytes = 1
	} else {
		numBytes = (length + 7) / 8 // Calculate the number of bytes needed to represent the given number of bits
	}

	bits := make([]byte, numBytes)

	_, err := rand.Read(bits)
	if err != nil {
		log.Fatal(err)
	}

	// Set remaining bits to 0 if length is not multiple of 8
	remainingBits := length % 8
	if remainingBits != 0 {
		bits[numBytes-1] &= (1 << remainingBits) - 1
	}
	return bits
}

// AllOnesBytes generates a slice of bytes of a given length, all set to 1
func AllOnesBytes(length int) []byte {
	b := make([]byte, length)
	for i := range b {
		b[i] = 1
	}
	return b
}

// RandomSelectionBits generates a slice of m random selection bits (0 or 1)
func RandomSelectionBits(m int) []byte {
	// Create a new random source and random generator
	src := mathRand.NewSource(time.Now().UnixNano())
	rnd := mathRand.New(src)

	SelectionBits := make([]byte, m)
	for i := range SelectionBits {
		SelectionBits[i] = byte(rnd.Intn(2)) // Generate a random byte 0 or 1
	}
	return SelectionBits
}

// divideMatrix divides a matrix of size rows x cols into smaller matrices of size rows x rows.
// The last matrix is padded with zeros if necessary to create a square matrix of size rows x rows as the other matrices.
func divideMatrix(matrix [][]byte, rows int, cols int) [][][]byte {

	if rows > cols {
		panic("divideMatrix: cols must be greater than or equal to rows")
	}

	// Calculate the number of rows x rows matrices
	numBaseMatrices := cols / rows
	numMatrices := numBaseMatrices

	// Add one matrix if there is padding
	remaining_cols := cols % rows
	padding_cols := 0
	if remaining_cols != 0 {
		padding_cols = rows - remaining_cols
	}
	if padding_cols > 0 {
		numMatrices += 1
	}

	// Initialize the result slice with size numMatrices
	result := make([][][]byte, numMatrices)

	for i := 0; i < numBaseMatrices; i++ {
		// Initialize smallMatrix for each sub-matrix of size rows x rows
		smallMatrix := make([][]byte, 0, len(matrix))

		for _, row := range matrix {

			start := i * rows
			end := start + rows

			// Append the row slice directly to the smallMatrix
			smallMatrix = append(smallMatrix, row[start:end])

		}
		// Append the smallMatrix to the result
		result = append(result, smallMatrix)
	}

	// If there is padding, add the last matrix
	if padding_cols > 0 {
		// Initialize smallMatrix for the last sub-matrix
		smallMatrix := make([][]byte, len(matrix))
		for _, row := range matrix {
			start := numBaseMatrices * rows
			end := start + (cols % rows) // Only take the non-padded elements
			paddedRow := make([]byte, rows)
			copy(paddedRow, row[start:end])
			// Padding will automatically be zeros in the remaining positions
			smallMatrix = append(smallMatrix, paddedRow)
		}
		result = append(result, smallMatrix)

	}
	return result
}

// EklundhTransposeInner transposes a matrix using Eklundh's algorithm.
// The matrix must be square and have a dimension that is a power of 2.
func EklundhTransposeInner(matrix [][]byte) [][]byte {
	dimension := len(matrix)

	if dimension == 1 {
		return matrix
	}

	// swapDimension is the dimension of the sub-matrix that is being swapped.
	// It starts at 1 and doubles each iteration until it reaches k. E.g. 1, 2, 4, 8, 16, ...
	swapDimension := 1 // Incremented by power of 2 each iteration

	for swapDimension < dimension {
		for i1 := swapDimension; i1 < dimension; i1 += swapDimension * 2 {
			for i2 := 0; i2 < swapDimension; i2++ {

				// // Index rows with values to be swapped
				topRow := matrix[i1+i2-swapDimension]
				bottomRow := matrix[i1+i2]

				for j1 := swapDimension; j1 < dimension; j1 += swapDimension * 2 {
					for j2 := 0; j2 < swapDimension; j2++ {

						topRow[j1+j2], bottomRow[j1+j2-swapDimension] = bottomRow[j1+j2-swapDimension], topRow[j1+j2]
					}
				}
				matrix[i1+i2-swapDimension] = topRow
				matrix[i1+i2] = bottomRow
			}
		}
		swapDimension *= 2
	}
	return matrix
}

// Function is responsible for dividing the matrix of m x k into smaller matrices of k x k.
// Each of the smaller matrices is then transposed using Eklundh's algorithm.
// The method required that k is a power of 2 else it fails.
// The method is multithreaded if multithreaded is set to true. This is done by using goroutines for each sub-matrix.
func EklundhTranspose(matrix [][]byte, multithreaded bool) [][]byte {

	rows := len(matrix)
	cols := len(matrix[0])

	matrices := divideMatrix(matrix, rows, cols) // Divide the matrix into smaller matrices of size rows x rows (k x k)

	if multithreaded {
		var wg sync.WaitGroup
		wg.Add(len(matrices)) // goroutines to wait for in waitgroup

		for i, mat := range matrices {
			go func(i int, mat [][]byte) {
				defer wg.Done()
				matrices[i] = EklundhTransposeInner(mat) // Transpose each sub-matrix concurrently on a separate goroutine
			}(i, mat) // Pass i and mat as arguments to the anonymous function to avoid race conditions on i
		}
		wg.Wait() // Wait for all goroutines to complete

	} else {
		// Transpose each sub-matrix sequentially on the main thread if multithreaded is false
		for i, mat := range matrices {
			matrices[i] = EklundhTransposeInner(mat)
		}
	}

	// Initialize the final transposed matrix
	transposed := make([][]byte, cols) // cols rows
	for i := range transposed {
		transposed[i] = make([]byte, rows) // rows columns
	}

	// If padding, calculate amount of rows that need to be unpadded (not included in the transposed matrix)
	padding_rows := 0
	if cols%rows != 0 {
		padding_rows = rows - (cols % rows)
	}

	// Stack the transposed matrices vertically
	iterations := len(matrices)*len(matrices[0]) - padding_rows // Exclude the padding rows
	currentRow := 0
	for _, mat := range matrices {
		for _, row := range mat {
			transposed[currentRow] = row // Add each row of the sub-matrices to the transposed matrix
			currentRow++

			// Break the loops before the padded row is reached in mat to avoid index out of range in transposed
			if currentRow == iterations {
				break
			}
		}
	}

	return transposed
}

// TransposeMatrix transposes a matrix using Eklundh's algorithm.
// The matrix must be square and have a dimension that is a power of 2.
// NOTICE: THIS METHOD IS NOT USED IN THE OT EXTENSION PROTOCOL.
// IT WAS AN EARLIER VERSION OF THE EKLUNDH TRANSPOSE METHOD USING RECURSION,
// WHICH WAS TERRIBLE AND THEREFORE REPLACED BY THE ITERATIVE VERSION.
func eklundhTransposeMatrixRecursiveInner(matrix [][]byte) [][]byte {

	rows := len(matrix)
	cols := len(matrix[0])

	// Padding if necessary to make the matrix square
	maxSize := max(rows, cols)

	// Base case for recursion
	if maxSize == 1 {
		return matrix
	}

	// Dividing the padded matrix into four sub-matrices
	rowMid := maxSize / 2
	colMid := maxSize / 2

	//var wg sync.WaitGroup

	A := makeSubMatrix(matrix, 0, rowMid, 0, colMid)             // Top left
	B := makeSubMatrix(matrix, 0, rowMid, colMid, maxSize)       // Top right
	C := makeSubMatrix(matrix, rowMid, maxSize, 0, colMid)       // Bottom left
	D := makeSubMatrix(matrix, rowMid, maxSize, colMid, maxSize) // Bottom right

	// Recursively transpose sub-matrices
	A = eklundhTransposeMatrixRecursiveInner(A)
	B = eklundhTransposeMatrixRecursiveInner(B)
	C = eklundhTransposeMatrixRecursiveInner(C)
	D = eklundhTransposeMatrixRecursiveInner(D)

	// Recursively transpose sub-matrices concurrently
	// wg.Add(4)
	// go func() {
	// 	defer wg.Done()
	// 	A = eklundhTransposeMatrix(A)
	// }()
	// go func() {
	// 	defer wg.Done()
	// 	B = eklundhTransposeMatrix(B)
	// }()
	// go func() {
	// 	defer wg.Done()
	// 	C = eklundhTransposeMatrix(C)
	// }()
	// go func() {
	// 	defer wg.Done()
	// 	D = eklundhTransposeMatrix(D)
	// }()
	// wg.Wait()
	return mergeSubMatrices(A, B, C, D)
}

// makeSubMatrix creates a sub-matrix from the given indexes of a matrix.
// NOTICE: THIS METHOD IS NOT USED IN THE OT EXTENSION PROTOCOL.
// IT WAS AN EARLIER VERSION OF THE EKLUNDH TRANSPOSE METHOD USING RECURSION,
// WHICH WAS TERRIBLE AND THEREFORE REPLACED BY THE ITERATIVE VERSION.
func makeSubMatrix(matrix [][]byte, rowStart, rowEnd, colStart, colEnd int) [][]byte {

	subMatrix := make([][]byte, rowEnd-rowStart)
	for i := range subMatrix {
		subMatrix[i] = matrix[rowStart+i][colStart:colEnd]
	}
	return subMatrix
}

// mergeSubMatrices merges four sub-matrices into a single matrix
// NOTICE: THIS METHOD IS NOT USED IN THE OT EXTENSION PROTOCOL.
// IT WAS AN EARLIER VERSION OF THE EKLUNDH TRANSPOSE METHOD USING RECURSION,
// WHICH WAS TERRIBLE AND THEREFORE REPLACED BY THE ITERATIVE VERSION.
func mergeSubMatrices(A, B, C, D [][]byte) [][]byte {

	topSideRows := len(A)      // Same as rowsA, rowsB
	leftSideCols := len(A[0])  // Same as colsA, colsC
	bottomSideRows := len(C)   // Same as rowsC, rowsD
	rightSideCols := len(B[0]) // Same as colsB, colsD

	totalRows := topSideRows + bottomSideRows
	totalCols := leftSideCols + rightSideCols

	matrix := make([][]byte, totalRows)
	for i := range matrix {
		matrix[i] = make([]byte, totalCols)
	}

	// Merge sub-matrices (See figure 2 in ALSZ paper)
	for i := 0; i < totalRows; i++ {
		for j := 0; j < totalCols; j++ {
			if i < topSideRows && j < leftSideCols {
				// Top left quadrant is A
				matrix[i][j] = A[i][j]
			} else if i < topSideRows && j >= leftSideCols {
				// Top right quadrant is swapped to C
				matrix[i][j] = C[i][j-leftSideCols]
			} else if i >= topSideRows && j < leftSideCols {
				// Bottom left quadrant is swapped to B
				matrix[i][j] = B[i-topSideRows][j]
			} else {
				// Bottom right quadrant is D
				matrix[i][j] = D[i-topSideRows][j-leftSideCols]
			}
		}
	}

	return matrix
}

// max returns the larger of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func TestDivideMatrix() {
	matrix := [][]byte{
		{0, 1, 0, 1, 1, 0, 0, 1, 1},
		{0, 1, 0, 0, 0, 0, 0, 1, 1},
	}

	rows := len(matrix)
	cols := len(matrix[0])

	for _, row := range matrix {
		fmt.Println(row)
	}
	fmt.Println()

	dividedMatrices := divideMatrix(matrix, rows, cols)

	for _, mat := range dividedMatrices {
		for _, row := range mat {
			fmt.Println(row)
		}
		fmt.Println()
	}

}

func TestEklundhTranspose() {
	matrix := [][]byte{
		{1, 0, 1, 0},
		{1, 1, 0, 0},
	}

	fmt.Println("Original Matrix:")
	for _, row := range matrix {
		fmt.Println(row)
	}

	transposedMatrix := EklundhTranspose(matrix, false)

	fmt.Println("\nTransposed Matrix:")
	for _, row := range transposedMatrix {
		fmt.Println(row)
	}

}

// Regular matric transpose using a naive algorithm of O(n^2)
// Used in the OTExtensionProtocolTranspose for transposing the matrices T and Q.
func TransposeMatrix(matrix [][]byte) [][]byte {
	newMatrix := make([][]byte, len(matrix[0]))

	for i := 0; i < len(matrix[0]); i++ {
		newMatrix[i] = make([]byte, len(matrix))
		for j := 0; j < len(matrix); j++ {
			newMatrix[i][j] = matrix[j][i]
		}
	}
	return newMatrix
}

// printMatrices prints a slice of matrices of type [][][]byte.
// Used for debugging.
func printMatrices(matrices [][][]byte) {
	for _, matrix := range matrices {
		for _, row := range matrix {
			for _, val := range row {
				fmt.Printf("%d ", val)
			}
			fmt.Println()
		}
		fmt.Println()
	}
}
