package benchmark

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

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

				// Index rows with values to be swapped
				//topRow := matrix[i1+i2]
				//bottomRow := matrix[i1+i2+swapDimension]

				// OUTCOMMENT TO SEE THE SWAPPING ORDER
				println("")
				fmt.Println("topRow: ", strconv.Itoa(i1+i2), "bottomRow: ", strconv.Itoa(i1+i2+swapDimension))

				for j1 := swapDimension; j1 < dimension; j1 += swapDimension * 2 {
					for j2 := 0; j2 < swapDimension; j2++ {

						// OUTCOMMENT TO SEE THE SWAPPING ORDER
						fmt.Println("Currently swapping: leftIndex: [][]", strconv.Itoa(i1+i2+swapDimension), strconv.Itoa(j1+j2), "rightIndex: [][] ",
							strconv.Itoa(i1+i2), strconv.Itoa(j1+j2+swapDimension))

					}
				}
			}
		}
		// OUTCOMMENT TO SEE THE SWAPPING ORDER
		//println("")
		swapDimension *= 2
	}

	return matrix
}

func EklundhTransposeInnerSøren(matrix [][]byte) [][]byte {
	dimension := len(matrix)

	if dimension == 1 {
		return matrix
	}

	// swapDimension is the dimension of the sub-matrix that is being swapped.
	// It starts at 1 and doubles each iteration until it reaches k. E.g. 1, 2, 4, 8, 16, ...
	swapDimension := 1 // Incremented by power of 2 each iteration
	for swapDimension < dimension {

	Yeehaw:
		for i := 0; i < dimension; i++ {
			var j int
			if i < swapDimension {
				j = i
			} else {
				j = i + (i * swapDimension)
			}

			fmt.Println("I:", i)
			fmt.Println("J:", j)
			if dimension < (j+swapDimension) || dimension < j {
				fmt.Println("J too large and breaks", j)
				break Yeehaw
			}

			topRow := make([]byte, swapDimension)
			bottomRow := make([]byte, swapDimension)
			copy(topRow, matrix[j])                  // Row of the top sub-matrix currently being swapped
			copy(bottomRow, matrix[j+swapDimension]) // Row of the bottom sub-matrix currently being swapped

			fmt.Println(topRow)
			fmt.Println(bottomRow)
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
	startTime := time.Now()

	rows := len(matrix)
	cols := len(matrix[0])

	// Timing division
	divisionStart := time.Now()
	matrices := divideMatrix(matrix, rows, cols) // Divide the matrix
	divisionTime := time.Since(divisionStart)

	// Timing the transposition
	transposeStart := time.Now()
	if multithreaded {
		var wg sync.WaitGroup
		wg.Add(len(matrices))

		for i, mat := range matrices {
			go func(i int, mat [][]byte) {
				defer wg.Done()
				matrices[i] = EklundhTransposeInner(mat)
			}(i, mat)
		}
		wg.Wait()
	} else {
		for i, mat := range matrices {
			matrices[i] = EklundhTransposeInner(mat)
		}
	}
	transposeTime := time.Since(transposeStart)

	// Timing final assembly
	assemblyStart := time.Now()
	transposed := make([][]byte, cols)
	for i := range transposed {
		transposed[i] = make([]byte, rows)
	}

	padding_rows := 0
	if cols%rows != 0 {
		padding_rows = rows - (cols % rows)
	}

	iterations := len(matrices)*len(matrices[0]) - padding_rows
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
	assemblyTime := time.Since(assemblyStart)

	totalTime := time.Since(startTime)

	fmt.Printf("Division Time: %v\n", divisionTime)
	fmt.Printf("Transpose Time: %v\n", transposeTime)
	fmt.Printf("Assembly Time: %v\n", assemblyTime)
	fmt.Printf("Total Time: %v\n", totalTime)

	return transposed
}
