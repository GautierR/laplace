package conduction

import (
	"conduction/utils"
	"fmt"
)

type Element struct {
	*Geometry
	*State
	prevState *State
	eNum      int
	xPosition float64
	start     float64
	end       float64
	length    float64
}

type Geometry struct {
	Source   HeatSource
	Material Material
}

// Compute the element length (used for simulation definition)
func (e *Element) setLength() {
	e.length = e.end - e.start
}

// NewEquidistantGrid return an equidistant grid of element pointers.
// The first and last element of the grid has a length of delta / 2.
func NewEquidistantGrid(start float64, end float64, n int) (grid []*Element, err error) {
	delta := (end - start) / float64(n-1)

	// Allocate the grid
	grid = make([]*Element, n)

	// Iterate over the whole domain
	currentPosition := start

	// First element
	grid[0] = &Element{
		eNum:      0,
		xPosition: currentPosition,
		start:     currentPosition,
		end:       currentPosition + delta/2,
		State:     &State{T: 0},
		prevState: &State{T: 0},
	}
	grid[0].setLength()
	currentPosition += grid[0].length

	for i := 1; i < n-1; i++ {
		grid[i] = &Element{
			eNum:      i,
			xPosition: currentPosition + delta/2,
			start:     currentPosition,
			end:       currentPosition + delta,
			State:     &State{T: 0},
			prevState: &State{T: 0},
		}
		grid[i].setLength()
		currentPosition += grid[i].length

	}

	// Last element
	grid[n-1] = &Element{
		eNum:      n - 1,
		xPosition: currentPosition + delta/2,
		start:     currentPosition,
		end:       currentPosition + delta/2,
		State:     &State{T: 0},
		prevState: &State{T: 0},
	}
	grid[n-1].setLength()
	currentPosition += grid[n-1].length

	if !utils.AlmostEqual(currentPosition, end) {
		err = fmt.Errorf("error in NewEquidistantGrid, last position [%v] does not match"+
			" the end of the domain [%v]", currentPosition, end)
	}

	return
}

func (e *Element) DistanceFromElement(start *Element) (dist float64) {
	return e.xPosition - start.xPosition
}

func (s *Simulation) ElementAt(idx int) *Element {
	return s.elements[idx]
}

func (s *Simulation) PositionAt(idx int) float64 {
	return s.elements[idx].xPosition
}

func (s *Simulation) TemperatureAt(idx int) float64 {
	return s.elements[idx].T
}

func (s *Simulation) GetXData() []float64 {
	var xData []float64
	for _, elem := range s.elements {
		xData = append(xData, elem.xPosition)
	}
	return xData
}

func (s *Simulation) GetTData() []float64 {
	var tData []float64
	for _, elem := range s.elements {
		tData = append(tData, elem.T)
	}
	return tData
}
