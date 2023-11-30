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
	"strconv"
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

// Struct to store ciphertexts the two messages M0 and M1
type CiphertextPair struct {
	Ciphertext0 *elgamal.Ciphertext
	Ciphertext1 *elgamal.Ciphertext
}

type Seed struct {
	Seed0 *big.Int
	Seed1 *big.Int
}

type ByteCiphertextPair struct {
	Y0 []byte
	Y1 []byte
}

func PseudoRandomGenerator(seed *big.Int, bitLength int) ([]uint8, error) {
	if bitLength <= 0 {
		return nil, fmt.Errorf("bitLength must be positive")
	}
	output := make([]uint8, 0, bitLength) // Allocate space for the array

	// Convert seed to a byte slice
	seedBytes := seed.Bytes()

	for len(output) < bitLength {
		hash := sha256.Sum256(seedBytes)
		for _, b := range hash[:] {
			for i := 0; i < 8 && len(output) < bitLength; i++ {
				bit := (b >> (7 - i)) & 1 // Extract each bit
				output = append(output, uint8(bit))
			}
		}

		// Increment the seed
		seed = new(big.Int).Add(seed, big.NewInt(1))
		seedBytes = seed.Bytes()
	}

	return output, nil
}

// Hash creates a hash of the input data with a specified byte length.
func Hash(originalData []byte, byteLength int) []byte {

	// Create a copy of the original data to avoid modifying it
	data := make([]byte, len(originalData))
	copy(data, originalData)

	fullHash := make([]byte, 0, byteLength)

	// SHA-256 produces a hash of 32 bytes (256 bits)
	// We keep hashing and concatenating until we reach the desired byte length
	for len(fullHash) < byteLength {
		hash := sha256.Sum256(data)
		fullHash = append(fullHash, hash[:]...)

		// Modify data slightly for the next iteration to produce a different hash
		// For example, append a byte that represents the current length of fullHash
		data = append(data, byte(len(fullHash)))
	}

	// Truncate the hash to the exact byte length required
	return fullHash[:byteLength]
}

// PrintBinaryString prints the binary representation of a []byte slice.
func PrintBinaryString(bytes []byte) {
	binaryString := ""
	for _, b := range bytes {
		binaryString += fmt.Sprintf("%08b", b) // Convert each byte to an 8-bit binary string
	}
	fmt.Println("As binary string:", binaryString)
}

// XOR takes a variable number of arguments (strings and ints),
// performs a bitwise XOR operation on all of them, and returns the result as a string.
func XOR(args ...interface{}) (string, error) {
	var xorResult int

	for _, arg := range args {
		switch v := arg.(type) {
		case string:
			val, err := strconv.Atoi(v)
			if err != nil {
				return "", err
			}
			xorResult ^= val
		case int:
			xorResult ^= v
		case uint8:
			xorResult ^= int(v)
		default:
			return "", fmt.Errorf("unsupported type")
		}
	}

	return strconv.Itoa(xorResult), nil
}

// RandomBytes generates a slice of random bytes of a given length
func RandomBytes(length int) []byte {
	b := make([]byte, length)
	_, err := rand.Read(b)
	// Handle the error here. In production code, you might want to pass it up the call stack
	if err != nil {
		log.Fatal(err)
	}
	return b
}

// AllOnesBytes generates a slice of bytes of a given length, all set to 1
func AllOnesBytes(length int) []byte {
	b := make([]byte, length)
	for i := range b {
		b[i] = 1
	}
	return b
}

// AllTwosBytes generates a slice of bytes of a given length, all set to 2
func AllTwosBytes(length int) []byte {
	b := make([]byte, length)
	for i := range b {
		b[i] = 2
	}
	return b
}

// RandomSelectionBits generates a slice of m random selection bits (0 or 1)
func RandomSelectionBits(m int) []uint8 {
	// Create a new random source and random generator
	src := mathRand.NewSource(time.Now().UnixNano())
	rnd := mathRand.New(src)

	bits := make([]uint8, m)
	for i := range bits {
		bits[i] = uint8(rnd.Intn(2)) // Generates a random integer 0 or 1
	}
	return bits
}

func divideMatrix(matrix [][]uint8, rows int, cols int) [][][]uint8 {

	if rows > cols {
		panic("divideMatrix: cols must be greater than or equal to rows")
	}

	// Calculate the number of rows x rows matrices

	numBaseMatrices := cols / rows
	remaining_cols := cols % rows
	padding_cols := 0
	if remaining_cols != 0 {
		padding_cols = rows - remaining_cols
	}
	numMatrices := numBaseMatrices // Add one matrix if there is padding
	if padding_cols > 0 {
		numMatrices += 1
	}

	// Initialize the result slice with size numMatrices
	result := make([][][]uint8, numMatrices)

	for i := 0; i < numBaseMatrices; i++ {
		// Initialize smallMatrix for each sub-matrix
		smallMatrix := make([][]uint8, 0, len(matrix))

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
		smallMatrix := make([][]uint8, len(matrix))
		for _, row := range matrix {
			start := numBaseMatrices * rows
			end := start + (cols % rows) // Only take the non-padded elements
			paddedRow := make([]uint8, rows)
			copy(paddedRow, row[start:end])
			// Padding will automatically be zeros in the remaining positions
			smallMatrix = append(smallMatrix, paddedRow)
		}
		result = append(result, smallMatrix)

	}
	return result
}

func EklundhTransposeInner(matrix [][]uint8) [][]uint8 {

	dimension := len(matrix)

	if dimension == 1 {
		return matrix
	}

	// swapDimension is the dimension of the sub-matrix that is being swapped.
	// It starts at 1 and doubles each iteration until it reaches k. E.g. 1, 2, 4, 8, 16, ...
	swapDimension := 1

	for swapDimension < dimension {

		for i := 0; i < dimension; i += 2 * swapDimension { // number of sub-matrices to swap in each iteration is k/(2*swapDimension)

			for j := 0; j < swapDimension; j++ {
				for l := 0; l < swapDimension; l++ {

					matrix[i+j][i+swapDimension+l], matrix[i+swapDimension+l][i+j] = matrix[i+swapDimension+l][i+j], matrix[i+j][i+swapDimension+l]

				}
			}
		}
		swapDimension *= 2
	}

	return matrix

}

// Function is responsible for dividing the matrix of m x k into smaller matrices of k x k. Also pads the last matrix if necessary.
func EklundhTranspose(matrix [][]uint8, multithreaded bool) [][]uint8 {

	rows := len(matrix)    // number of rows
	cols := len(matrix[0]) // number of columns

	matrices := divideMatrix(matrix, rows, cols)

	if multithreaded {
		var wg sync.WaitGroup
		wg.Add(len(matrices)) // goroutines to wait for in waitgroup

		for i, mat := range matrices {
			go func(i int, mat [][]uint8) {
				defer wg.Done()
				matrices[i] = EklundhTransposeInner(mat)
			}(i, mat) // Pass i and mat as arguments to the anonymous function to avoid race conditions on i
		}

		wg.Wait() // Wait for all goroutines to complete
	} else {
		for i, mat := range matrices {
			matrices[i] = EklundhTransposeInner(mat)
		}
	}

	// Initialize the final transposed matrix
	transposed := make([][]uint8, cols) // cols rows
	for i := range transposed {
		transposed[i] = make([]uint8, rows) // rows columns
	}

	padding_rows := 0
	if cols%rows != 0 {
		padding_rows = rows - (cols % rows)
	}

	// Stack the transposed matrices vertically
	iterations := len(matrices)*len(matrices[0]) - padding_rows // Exclude the padding rows
	currentRow := 0
	for _, mat := range matrices {
		for _, row := range mat {
			transposed[currentRow] = row
			currentRow++
			if currentRow == iterations {
				break
			}
		}
	}

	// print("Transposed matrices: " + "\n")
	// printMatrices(matrices)

	// print("Transposed matrix: " + "\n")
	// for _, row := range transposed {
	// 	fmt.Println(row)
	// }

	return transposed
}

// TransposeMatrix transposes a matrix using Eklundh's algorithm
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
// Was used for the recursive version of Eklundh's algorithm.
func makeSubMatrix(matrix [][]byte, rowStart, rowEnd, colStart, colEnd int) [][]byte {

	subMatrix := make([][]byte, rowEnd-rowStart)
	for i := range subMatrix {
		subMatrix[i] = matrix[rowStart+i][colStart:colEnd]
	}
	return subMatrix
}

// mergeSubMatrices merges four sub-matrices into a single matrix
// Was used for the recursive version of Eklundh's algorithm.
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

// regular matric transpose.
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

// padMatrix pads the given matrix with zeros to make it square
// Method was used for the recursive version of Eklundh's algorithm.
func padMatrix(matrix [][]byte, size int) [][]byte {
	paddedMatrix := make([][]byte, size)
	for i := range paddedMatrix {
		paddedMatrix[i] = make([]byte, size)
		if i < len(matrix) {
			//for j, val := range matrix[i] {
			//paddedMatrix[i][j] = val
			copy(paddedMatrix[i], matrix[i])
			//}
		}
	}
	return paddedMatrix
}

// Method was used for the recursive version of Eklundh's algorithm.
func unpadMatrix(matrix [][]byte, rows, cols int) [][]byte {
	unpaddedMatrix := make([][]byte, rows)
	for i := range unpaddedMatrix {
		unpaddedMatrix[i] = matrix[i][:cols]
	}
	return unpaddedMatrix
}

// printMatrices prints a slice of matrices of type [][][]byte.
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

func PrintMatrix(matrix [][]uint8) {
	for _, row := range matrix {
		fmt.Println(row)
	}
	fmt.Println()
}

func DeepCopyMatrix(matrix [][]uint8) [][]uint8 {
	numRows := len(matrix)
	if numRows == 0 {
		return nil // or return [][]uint8{} if you prefer an empty matrix over nil
	}

	numCols := len(matrix[0])
	newMatrix := make([][]uint8, numRows)
	for i := range matrix {
		newMatrix[i] = make([]uint8, numCols)
		copy(newMatrix[i], matrix[i])
	}

	return newMatrix
}
