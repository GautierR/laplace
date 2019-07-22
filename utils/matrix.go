package utils

import (
	"fmt"
	"gonum.org/v1/gonum/mat"
)

func PrintDense(m *mat.BandDense) {
	r, c := m.Dims()
	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			fmt.Printf("     %v", m.At(i, j))
		}
		fmt.Printf("\n")
	}
	fmt.Printf("\n")
}
