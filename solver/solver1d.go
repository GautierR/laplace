package solver

import "gonum.org/v1/gonum/mat"

type Solver1D struct {
	A *mat.BandDense
	B *mat.VecDense
	T *mat.VecDense
}

func NewSolver1D(nE int) *Solver1D {
	return &Solver1D{
		A: mat.NewBandDense(nE, nE, 1, 1, nil),
		B: mat.NewVecDense(nE, nil),
		T: mat.NewVecDense(nE, nil),
	}
}

func (s *Solver1D) SetAMatrix(m *mat.BandDense) {
	s.A = m
}

func (s *Solver1D) SetBVector(v *mat.VecDense) {
	s.B = v
}

func (s *Solver1D) Solve() (err error) {
	err = s.T.SolveVec(s.A, s.B)
	return err
}
