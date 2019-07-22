package conduction

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type SimulationInput struct {
	Domains            []Domain            `json:"domain"`
	BoundaryConditions []BoundaryCondition `json:"boundary_conditions"`
	InitialConditions  InitialConditions   `json:"initial_conditions"`
	Parameters         Parameters          `json:"parameters"`
}

type Domain struct {
	Start    float64    `json:"start"`
	End      float64    `json:"end"`
	Source   HeatSource `json:"source"`
	Material Material   `json:"material"`
}

type HeatSource struct {
	Constant float64 `json:"constant"`
	Linear   float64 `json:"linear"`
}

type Material struct {
	Name         string  `json:"name"`
	Density      float64 `json:"density"`
	Conductivity float64 `json:"conductivity"`
	SpecificHeat float64 `json:"specific_heat"`
}

type BoundaryCondition struct {
	Start  float64 `json:"start"`
	End    float64 `json:"end"`
	BCType string  `json:"bc_type"`
	Value  float64 `json:"value"`
	Value2 float64 `json:"value_2"`
}

type InitialConditions struct {
	Temperature float64 `json:"temperature"`
}

type Parameters struct {
	IterationMax int     `json:"iteration_max"`
	Tolerance    float64 `json:"tolerance"`
	NElement     int     `json:"n_element"`
}

func NewSimulationInput(fileName string) *SimulationInput {
	var simulationInput SimulationInput

	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Panicf("Error while reading case file: %v \n", err)
	}

	err = json.Unmarshal([]byte(file), &simulationInput)
	if err != nil {
		log.Panicf("Error while loading case file: %v \n", err)
	}

	return &simulationInput
}

func (sI *SimulationInput) GetDomainLimits() (start float64, end float64) {
	start = sI.Domains[0].Start
	end = sI.Domains[len(sI.Domains)-1].End
	return
}

func (s *Simulation) ApplyGeometry(domains []Domain) {
	for _, domain := range domains {
		geometry := &Geometry{
			Source:   domain.Source,
			Material: domain.Material,
		}
		for _, element := range s.elements {
			if domain.Start <= element.xPosition &&
				element.xPosition <= domain.End+1e9 {
				element.Geometry = geometry
			}
		}
	}
}
