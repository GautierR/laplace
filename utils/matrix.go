package utils

import (
	"gonum.org/v1/gonum/mat"
	"math"
)

// PrintDenseMatrix
//func PrintDenseMatrix(m *mat.Dense) {
//	r, c := m.Dims()
//	for i := 0; i < r; i++ {
//		for j := 0; j < c; j++ {
//			fmt.Printf("     %v", m.At(i, j))
//		}
//		fmt.Printf("\n")
//	}
//	fmt.Printf("\n")
//}

//VecDiffNormL2 returns the L2-norm of a-b
func VecDiffNormL2(a, b *mat.VecDense) (normL2 float64) {
	for i := 0; i < a.Len(); i++ {
		normL2 += math.Pow(a.AtVec(i)-b.AtVec(i), 2)
	}
	normL2 = math.Sqrt(normL2)
	return
}
