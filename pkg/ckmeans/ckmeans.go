package ckmeans

import (
	"errors"
	"sort"
)

type F64Matrix [][]float64
type IntMatrix [][]int

func max(a int, b int) int {
	if b > a {
		return b
	}
	return a
}

func min(a int, b int) int {
	if b < a {
		return b
	}
	return a
}

func ssq(j int, i int, sums []float64, sumsOfSquares []float64) float64 {
	var sji float64 // s(j, i)
	if j > 0 {
		muji := (sums[i] - sums[j-1]) / float64(i-j+1) // mu(j, i)
		sji =
			sumsOfSquares[i] - sumsOfSquares[j-1] - float64(i-j+1)*muji*muji
	} else {
		sji = sumsOfSquares[i] - (sums[i]*sums[i])/float64(i+1)
	}
	if sji < 0 {
		return 0
	}
	return sji
}

func fillMatrixColumn(
	iMin int,
	iMax int,
	cluster int,
	matrix F64Matrix,
	backtrackMatrix IntMatrix,
	sums []float64,
	sumsOfSquares []float64) {
	if iMin > iMax {
		return
	}

	// Start at midpoint between iMin and iMax
	i := (iMin + iMax) / 2

	matrix[cluster][i] = matrix[cluster-1][i-1]
	backtrackMatrix[cluster][i] = i

	jlow := cluster // the lower end for j

	if iMin > cluster {
		jlow = max(jlow, backtrackMatrix[cluster][iMin-1])
	}
	jlow = max(jlow, backtrackMatrix[cluster-1][i])

	jhigh := i - 1 // the upper end for j
	if iMax < len(matrix[0])-1 {
		jhigh = min(jhigh, backtrackMatrix[cluster][iMax+1])
	}

	for j := jhigh; j >= jlow; j-- {
		sji := ssq(j, i, sums, sumsOfSquares)

		if sji+matrix[cluster-1][jlow-1] >= matrix[cluster][i] {
			break
		}

		// Examine the lower bound of the cluster border
		sjlowi := ssq(jlow, i, sums, sumsOfSquares)

		ssqjlow := sjlowi + matrix[cluster-1][jlow-1]

		if ssqjlow < matrix[cluster][i] {
			// Shrink the lower bound
			matrix[cluster][i] = ssqjlow
			backtrackMatrix[cluster][i] = jlow
		}
		jlow++

		ssqj := sji + matrix[cluster-1][j-1]
		if ssqj < matrix[cluster][i] {
			matrix[cluster][i] = ssqj
			backtrackMatrix[cluster][i] = j
		}
	}

	fillMatrixColumn(iMin, i-1, cluster, matrix, backtrackMatrix, sums, sumsOfSquares)
	fillMatrixColumn(i+1, iMax, cluster, matrix, backtrackMatrix, sums, sumsOfSquares)
}

func fillMatrices(data []float64, matrix F64Matrix, backtrackMatrix IntMatrix) {
	var nValues int = len(matrix[0])

	// Shift values by the median to improve numeric stability
	shift := data[nValues/2]

	// Cumulative sum and cumulative sum of squares for all values in data array
	sums := []float64{}
	sumsOfSquares := []float64{}

	// Initialize first column in matrix & backtrackMatrix
	for i := 0; i < nValues; i++ {
		shiftedValue := data[i] - shift
		if i == 0 {
			sums = append(sums, shiftedValue)
			sumsOfSquares = append(sumsOfSquares, shiftedValue*shiftedValue)
		} else {
			sums = append(sums, sums[i-1]+shiftedValue)
			sumsOfSquares = append(sumsOfSquares, sumsOfSquares[i-1]+shiftedValue*shiftedValue)
		}

		// Initialize for cluster = 0
		matrix[0][i] = ssq(0, i, sums, sumsOfSquares)
		backtrackMatrix[0][i] = 0
	}

	// Initialize the rest of the columns
	var iMin int
	for cluster := 1; cluster < len(matrix); cluster++ {
		if cluster < len(matrix)-1 {
			iMin = cluster
		} else {
			// No need to compute matrix[K-1][0] ... matrix[K-1][N-2]
			iMin = nValues - 1
		}

		fillMatrixColumn(iMin, nValues-1, cluster, matrix, backtrackMatrix, sums, sumsOfSquares)
	}
}

func allEqual(vals []float64) bool {
	for i := 1; i < len(vals); i++ {
		if vals[i] != vals[0] {
			return false
		}
	}
	return true
}

func makeF64Matrix(ncol, nrow int) F64Matrix {
	result := make(F64Matrix, ncol)
	for i := range result {
		result[i] = make([]float64, nrow)
	}
	return result
}

func makeIntMatrix(ncol, nrow int) IntMatrix {
	result := make(IntMatrix, ncol)
	for i := range result {
		result[i] = make([]int, nrow)
	}
	return result
}

func Ckmeans(x []float64, nClusters int) ([][]float64, error) {
	if nClusters > len(x) {
		return nil, errors.New("cannot generate more classes than there are data values")
	}

	sorted := make([]float64, len(x))
	copy(sorted, x)
	sort.Float64s(sorted)

	// if all of the input values are identical, there's one cluster
	// with all of the input in it.
	if allEqual(sorted) {
		result := make([][]float64, 1)
		result[0] = sorted
		return result, nil
	}

	// named 'S' originally
	matrix := makeF64Matrix(nClusters, len(sorted))
	// named 'J' originally
	backtrackMatrix := makeIntMatrix(nClusters, len(sorted))

	// This is a dynamic programming way to solve the problem of minimizing
	// within-cluster sum of squares. It's similar to linear regression
	// in this way, and this calculation incrementally computes the
	// sum of squares that are later read.
	fillMatrices(sorted, matrix, backtrackMatrix)

	// The real work of Ckmeans clustering happens in the matrix generation:
	// the generated matrices encode all possible clustering combinations, and
	// once they're generated we can solve for the best clustering groups
	// very quickly.
	clusters := make([][]float64, nClusters)
	clusterRight := len(backtrackMatrix[0]) - 1

	// Backtrack the clusters from the dynamic programming matrix. This
	// starts at the bottom-right corner of the matrix (if the top-left is 0, 0),
	// and moves the cluster target with the loop.
	for cluster := len(backtrackMatrix) - 1; cluster >= 0; cluster-- {
		clusterLeft := backtrackMatrix[cluster][clusterRight]

		// fill the cluster from the sorted input by taking a slice of the
		// array. the backtrack matrix makes this easy - it stores the
		// indexes where the cluster should start and end.
		clusters[cluster] = sorted[clusterLeft : clusterRight+1]

		if cluster > 0 {
			clusterRight = clusterLeft - 1
		}
	}

	return clusters, nil
}
