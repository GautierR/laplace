package conduction

import (
	"conduction/utils"
	"fmt"
	"math"
)

type Element struct {
	*Geometry
	*State            // State contain the physical properties of the element (i.e. T,..)
	prevState *State  // Previous state
	eNum      int     // Element number, correspond to the index in Simulation grid
	xPosition float64 // [m] Position of the node in the x axis
	start     float64 // [m] Starting point of the element
	end       float64 // [m] Ending point of the element
	length    float64 // [m] Element length

	aW *float64 // A Matrix coefficient at west
	aP *float64 // A Matrix coefficient at point
	aE *float64 // A Matrix coefficient at east
	b  *float64 // b source value

	deltaW                 float64 // [m] Distance between W and P
	deltaE                 float64 // [m] Distance between E and P
	interfaceConductivityW float64 // [W/m.K] West interface thermal conductivity
	interfaceConductivityE float64 // [W/m.K] East interface thermal conductivity

	previousElement *Element // Pointer to the previous element
	nextElement     *Element // Pointer to the next element
}

type Geometry struct {
	Source   HeatSource
	Material Material
}

// NewEquidistantGrid return an equidistant grid of element pointers.
// The first and last element of the grid has a length of delta / 2.
func NewEquidistantGrid(start float64, end float64, n int) (grid []*Element, err error) {
	delta := (end - start) / float64(n-1)

	// Allocate the grid
	grid = make([]*Element, n)

	// Iterate over the whole domain
	currentPosition := start

	var xPosition float64
	var currentStart float64
	var currentEnd float64

	for i := range grid {
		if i == 0 {
			// First element
			xPosition = currentPosition
			currentStart = currentPosition
			currentEnd = currentPosition + delta/2

		} else if i == n-1 {
			// Last element
			xPosition = currentPosition + delta/2
			currentStart = currentPosition
			currentEnd = currentPosition + delta/2

		} else {
			xPosition = currentPosition + delta/2
			currentStart = currentPosition
			currentEnd = currentPosition + delta
		}

		grid[i] = &Element{
			eNum:      i,
			xPosition: xPosition,
			start:     currentStart,
			end:       currentEnd,
			length:    currentEnd - currentStart,
			State:     NewState(),
			prevState: NewState(),
		}

		currentPosition += grid[i].length
	}

	if !utils.AlmostEqual(currentPosition, end) {
		err = fmt.Errorf("error in NewEquidistantGrid, last position [%v] does not match"+
			" the end of the domain [%v]", currentPosition, end)
	}

	return
}

func (e *Element) DistanceFromElement(start *Element) (dist float64) {
	if start != nil {
		return math.Abs(e.xPosition - start.xPosition)
	}
	return
}

func (e *Element) SetCoefficient() {
	if e.previousElement != nil {
		aW := e.interfaceConductivityW / e.deltaW
		*e.aW = -aW
		*e.aP += aW
	}
	if e.nextElement != nil {
		aE := e.interfaceConductivityE / e.deltaE
		*e.aE = -aE
		*e.aP += aE
	}

}

func (e *Element) SetSource() {
	*e.b = e.Source.Constant * e.length
}

func (e *Element) SetBoundaryCondition(condition BoundaryCondition) {
	condition.SetAtElement(e)
}
