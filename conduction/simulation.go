package conduction

import (
	"conduction/file"
	"conduction/solver"
	"conduction/utils"
	"fmt"
	"gonum.org/v1/gonum/mat"
	"log"
	"math"
)

type Simulation struct {
	solver1D           *solver.Solver1D
	elements           []*Element
	boundaryConditions []BoundaryCondition
	parameters         Parameters
	stateVector        *mat.VecDense
	prevStateVector    *mat.VecDense
}

func New1DConductionSimulation(simInput *SimulationInput) *Simulation {
	// Generate mesh
	nElement := simInput.Parameters.NElement

	// Get simulation domain limits
	start, end := simInput.GetDomainLimits()

	// Generate the mesh grid
	grid, err := NewEquidistantGrid(start, end, nElement)
	if err != nil {
		fmt.Printf("Grid generation error: %v \n", err)
	}

	// Generate simulation object
	sim := &Simulation{
		solver1D:           solver.NewSolver1D(nElement),
		elements:           grid,
		boundaryConditions: ParseBoundaryConditions(simInput.BoundaryConditions),
		parameters:         simInput.Parameters,
		prevStateVector:    mat.NewVecDense(nElement, nil),
	}

	// Link via pointer result vector with Simulation state
	sim.stateVector = sim.solver1D.X

	sim.ApplyGeometry(simInput.Domains)
	sim.InitializeElements(simInput)

	for _, elem := range sim.elements {
		elem.SetCoefficient()
		elem.SetSource()
	}
	sim.SetBoundaryConditions()

	return sim
}

func (s *Simulation) Start() {
	// Setup parameters
	normError := math.Inf(1)
	tolerance := s.parameters.Tolerance

	var iteration int
	iterationMax := s.parameters.IterationMax

	for normError > tolerance && iteration < iterationMax {
		normError, _ = s.Solve()
		iteration += 1
	}
	fmt.Printf("Iteration : %v", iteration)
	return
}

func (s *Simulation) Solve() (normError float64, err error) {
	err = s.solver1D.Solve()
	if err != nil {
		err = fmt.Errorf("Error in solver1D: %v \n", err)
	}
	normError = utils.VecDiffNormL2(s.stateVector, s.prevStateVector)
	s.UpdatePreviousState()

	return
}

func (s *Simulation) ExportData(inputFile string) {
	xData := s.GetXData()
	tData := s.GetTData()

	err := file.ExportToCsv(inputFile, xData, tData)
	if err != nil {
		log.Panicf("Error while exporting the results into csv: %v \n", err)
	}
}
