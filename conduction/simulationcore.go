package conduction

func (s *Simulation) ElementAt(idx int) *Element {
	return s.elements[idx]
}

func (s *Simulation) PositionAt(idx int) float64 {
	return s.elements[idx].xPosition
}

func (s *Simulation) NearestElementFrom(x float64) *Element {
	for _, elem := range s.elements {
		if elem.xPosition > x {
			if x-elem.previousElement.xPosition < elem.xPosition-x {
				return elem.previousElement
			}
			return elem
		}
	}
	return s.elements[len(s.elements)-1]
}

func (s *Simulation) TemperatureAt(idx int) float64 {
	return s.elements[idx].T()
}

func (s *Simulation) HeatFluxAt(idx int) float64 {
	return s.elements[idx].HeatFlux()
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
		tData = append(tData, elem.T())
	}
	return tData
}

func (s *Simulation) GetHFData() []float64 {
	var tData []float64
	for _, elem := range s.elements {
		tData = append(tData, elem.HeatFlux())
	}
	return tData
}

func (s *Simulation) LinkElements() {
	var previousElement *Element
	var nextElement *Element

	for i, elem := range s.elements {
		if i == 0 {
			previousElement = nil
			nextElement = s.elements[i+1]
		} else if i == len(s.elements)-1 {
			previousElement = s.elements[i-1]
			nextElement = nil
		} else {
			previousElement = s.elements[i-1]
			nextElement = s.elements[i+1]
		}

		elem.previousElement = previousElement
		elem.nextElement = nextElement
	}
}

func (s *Simulation) InitializeElements() {
	s.LinkElements()

	for _, elem := range s.elements {

		// Link element attributes to solver matrices
		A := s.solver1D.A
		B := s.solver1D.B
		state := s.stateVector
		prevState := s.prevStateVector
		elem.LinkToMatrices(A, B, state, prevState)

		// Set Initial Temperature
		elem.prevState.SetT(s.InitialConditions[elem.domainTag].Temperature)

		// Compute and update interface thermal conductivity
		elem.UpdateInterfaceConductivity(s.Parameters.InterfaceMethod)
	}
}

// SetBoundaryCondition
func (s *Simulation) SetBoundaryConditions() {
	for _, bC := range s.BoundaryConditions {
		tZ := s.ThermalZones[bC.Tag]
		if bC.IsAtPoint() {
			s.NearestElementFrom(bC.Start).SetBoundaryCondition(&bC, &tZ)
		} else {
			for _, elem := range s.elements {
				if elem.xPosition >= bC.Start && elem.xPosition <= bC.End {
					elem.SetBoundaryCondition(&bC, &tZ)
				}
			}
		}
	}
}

func (s *Simulation) UpdateState() {
	for _, elem := range s.elements {
		elem.prevState.SetT(elem.T())
		elem.SetHeatFlux()
	}
}
