package conduction

import (
	"log"
)

type BoundaryCondition interface {
	Start() float64
	End() float64
	SetAtElement(e *Element)
}

type TemperatureBoundaryCondition struct {
	BoundaryCondition
	start       float64
	end         float64
	Temperature float64
}

type HeatFluxBoundaryCondition struct {
	BoundaryCondition
	start    float64
	end      float64
	HeatFlux float64
}

type HeatTransferCoefficientBoundaryCondition struct {
	BoundaryCondition
	start                   float64
	end                     float64
	HeatTransferCoefficient float64
	InfinityTemperature     float64
}

func ParseBoundaryConditions(boundaryConditions []map[string]interface{}) (bCS []BoundaryCondition) {
	bCS = make([]BoundaryCondition, len(boundaryConditions))

	for i, bC := range boundaryConditions {
		if _, typeExist := bC["type"]; typeExist {
			switch bC["type"] {
			case "constant_temperature":
				bCS[i] = TemperatureBoundaryCondition{
					start:       bC["start"].(float64),
					end:         bC["end"].(float64),
					Temperature: bC["temperature"].(float64),
				}
			case "constant_heat_flux":
				bCS[i] = HeatFluxBoundaryCondition{
					start:    bC["start"].(float64),
					end:      bC["end"].(float64),
					HeatFlux: bC["heat_flux"].(float64),
				}
			case "heat_transfer_coefficient":
				bCS[i] = HeatTransferCoefficientBoundaryCondition{
					start:                   bC["start"].(float64),
					end:                     bC["end"].(float64),
					HeatTransferCoefficient: bC["heat_transfer_coefficient"].(float64),
					InfinityTemperature:     bC["infinity_temperature"].(float64),
				}
			}
		} else {
			log.Panicf("Missing boundary condition type.")
		}
	}
	return
}

func (bC TemperatureBoundaryCondition) Start() float64 {
	return bC.start
}

func (bC TemperatureBoundaryCondition) End() float64 {
	return bC.end
}

func (bC HeatFluxBoundaryCondition) Start() float64 {
	return bC.start
}

func (bC HeatFluxBoundaryCondition) End() float64 {
	return bC.end
}

func (bC HeatTransferCoefficientBoundaryCondition) Start() float64 {
	return bC.start
}

func (bC HeatTransferCoefficientBoundaryCondition) End() float64 {
	return bC.end
}

func (bC TemperatureBoundaryCondition) SetAtElement(e *Element) {
	if e.previousElement != nil {
		*e.aW = 0
	}
	if e.nextElement != nil {
		*e.aE = 0
	}
	*e.aP = 1
	*e.b = bC.Temperature
}

func (bC HeatFluxBoundaryCondition) SetAtElement(e *Element) {
	*e.b += bC.HeatFlux
}

func (bC HeatTransferCoefficientBoundaryCondition) SetAtElement(e *Element) {
	*e.aP += bC.HeatTransferCoefficient
	*e.b += bC.HeatTransferCoefficient * bC.InfinityTemperature
}
