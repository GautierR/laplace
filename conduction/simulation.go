package conduction

import (
	"conduction/file"
	"conduction/solver"
	"fmt"
	"log"
	"math"
)

type Simulation struct {
	solver1D           *solver.Solver1D
	elements           []*Element
	boundaryConditions []BoundaryCondition
	initialConditions  InitialConditions
	parameters         Parameters
}

func New1DConductionSimulation(simInput *SimulationInput) *Simulation {
	// Generate mesh
	nElement := simInput.Parameters.NElement
	start, end := simInput.GetDomainLimits()

	grid, err := NewEquidistantGrid(start, end, nElement)
	if err != nil {
		fmt.Printf("Grid generation error: %v \n", err)
	}

	sim := &Simulation{
		elements:           grid,
		boundaryConditions: simInput.BoundaryConditions,
		initialConditions:  simInput.InitialConditions,
		parameters:         simInput.Parameters,
	}
	sim.ApplyGeometry(simInput.Domains)
	sim.SetInitialTemperature(sim.initialConditions.Temperature)

	sim.solver1D = solver.NewSolver1D(nElement)

	sim.SetConductanceMatrix()
	sim.SetSourceVector()
	sim.SetBoundaryConditions()

	return sim
}

func (s *Simulation) Start() {
	// Setup parameters
	deltaT := math.Inf(1)
	tolerance := s.parameters.Tolerance

	var iteration int
	iterationMax := s.parameters.IterationMax

	// Temperature vector
	var err error

	for deltaT > tolerance && iteration < iterationMax {
		err = s.solver1D.Solve()
		if err != nil {
			log.Panicf("Error in solver1D: %v \n", err)
		}
		s.UpdateState()
		deltaT = s.GetTemperatureNorm()
		s.UpdatePreviousState()
		iteration += 1
	}

	return
}

func (s *Simulation) PrevElem(e *Element) *Element {
	return s.elements[e.eNum-1]
}

func (s *Simulation) NextElem(e *Element) *Element {
	return s.elements[e.eNum+1]
}

func (s *Simulation) SetInitialTemperature(temperature float64) {
	for _, elem := range s.elements {
		elem.prevState.T = temperature
	}
}

func (s *Simulation) UpdateState() {
	for idx, elem := range s.elements {
		elem.T = s.solver1D.T.AtVec(idx)
	}
}

func (s *Simulation) UpdatePreviousState() {
	for _, elem := range s.elements {
		elem.prevState = elem.State
	}
}

func (s *Simulation) GetTemperatureNorm() float64 {
	var norm float64
	for _, elem := range s.elements {
		norm += math.Pow(elem.T-elem.prevState.T, 2)
	}
	return math.Sqrt(norm)
}

func (s *Simulation) ExportData(inputFile string) {
	xData := s.GetXData()
	tData := s.GetTData()

	err := file.ExportToCsv(inputFile, xData, tData)
	if err != nil {
		log.Panicf("Error while exporting the results into csv: %v \n", err)
	}
}
