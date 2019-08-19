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
	Parameters         Parameters                  `json:"parameters"`
	Domains            map[string]Domain           `json:"domains"`
	ThermalZones       map[string]ThermalZone      `json:"thermal_zones"`
	Materials          map[string]Material         `json:"materials"`
	InitialConditions  map[string]InitialCondition `json:"initial_conditions"`
	BoundaryConditions []BoundaryCondition         `json:"boundary_conditions"`

	elements []*Element
	solver1D *solver.Solver1D

	stateVector     *mat.VecDense
	prevStateVector *mat.VecDense
}

type Parameters struct {
	Id              int     `json:"id"`
	IterationMax    int     `json:"iteration_max"`
	Tolerance       float64 `json:"tolerance"`
	NElement        int     `json:"n_element"`
	InterfaceMethod string  `json:"interface_method"`
}

type Domain struct {
	Start    float64  `json:"start"`
	End      float64  `json:"end"`
	Material string   `json:"material"`
	Geometry Geometry `json:"geometry"`
}

type ThermalZone struct {
	Type                    string  `json:"type"`
	FixedTemperature        float64 `json:"temperature"`               // [K] Local zone fixed temperature
	HeatFlow                float64 `json:"heat_flow"`                 // [W] Local heat applied on the zone
	HeatFlux                float64 `json:"heat_flux"`                 // [W/m²] Local heat flux applied on the zone
	InternalHeat            float64 `json:"internal_heat"`             // [W/m³] Internal heat generation
	HeatTransferCoefficient float64 `json:"heat_transfer_coefficient"` // [W/m².K] Heat transfer coefficient
	InfinityTemperature     float64 `json:"infinity_temperature"`      // [K] Temperature at infinity
}

type Material struct {
	Density      float64 `json:"density"`
	Conductivity float64 `json:"conductivity"`
	SpecificHeat float64 `json:"specific_heat"`
}

type InitialCondition struct {
	Temperature float64 `json:"temperature"`
}

func New1DConductionSimulation(fileName string) *Simulation {
	// Parse simulation file
	sim, err := ParseSimulationFromFile(fileName)
	if err != nil {
		log.Panicf("Error while parsing simulation : %v", err)
	}

	// Generate mesh
	nElement := sim.Parameters.NElement

	// Get simulation domain limits
	start, end := sim.GetDomainLimits()

	// Generate the mesh grid
	sim.elements, err = NewEquidistantGrid(start, end, nElement)
	if err != nil {
		fmt.Printf("Grid generation error: %v \n", err)
	}

	// Initialize the solver
	sim.solver1D = solver.NewSolver1D(nElement)

	// Create and link state vectors
	sim.prevStateVector = mat.NewVecDense(nElement, nil)
	sim.stateVector = sim.solver1D.X

	sim.ApplyGeometry()
	sim.InitializeElements()

	for _, elem := range sim.elements {
		elem.SetCoefficient()
	}

	sim.SetBoundaryConditions()

	return sim
}

func (s *Simulation) ApplyGeometry() {
	for name, domain := range s.Domains {
		material := s.Materials[domain.Material]
		for _, element := range s.elements {
			if domain.Start <= element.xPosition &&
				element.xPosition <= domain.End+1e-9 {
				element.domainTag = name
				element.Geometry = &domain.Geometry
				element.Material = &material
			}
		}
	}
}

func (s *Simulation) Start() {
	// Setup parameters
	normError := math.Inf(1)
	tolerance := s.Parameters.Tolerance

	var iteration int
	iterationMax := s.Parameters.IterationMax

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
	s.UpdateState()

	return
}

func (s *Simulation) ExportData(inputFile string) {
	xData := s.GetXData()
	tData := s.GetHFData()

	err := file.ExportToCsv(inputFile, xData, tData)
	if err != nil {
		log.Panicf("Error while exporting the results into csv: %v \n", err)
	}
}
