package conduction

// SetConductanceMatrix return []float64 array arranged to be able to create a BandDense matrix
func (s *Simulation) SetConductanceMatrix() {
	A := s.solver1D.A
	nElem := len(s.elements)

	var elem *Element
	var prevElem *Element
	var nextElem *Element

	// First element
	idx := 0
	elem = s.elements[idx]
	nextElem = s.NextElem(elem)
	conductivity := elem.Material.Conductivity
	deltaX := nextElem.DistanceFromElement(elem)

	aValue := 2 * conductivity / deltaX
	aNeighbour := conductivity / deltaX
	A.SetBand(idx, idx, aValue)
	A.SetBand(idx, idx+1, -aNeighbour)

	for idx = 1; idx < nElem-1; idx++ {
		elem = s.elements[idx]
		prevElem = s.PrevElem(elem)
		nextElem = s.NextElem(elem)
		conductivity = elem.Material.Conductivity
		deltaX = nextElem.DistanceFromElement(elem)

		aValue = 2 * conductivity / deltaX
		aNeighbour = conductivity / deltaX
		A.SetBand(idx, idx-1, -aNeighbour)
		A.SetBand(idx, idx, aValue)
		A.SetBand(idx, idx+1, -aNeighbour)
	}

	// Last element
	elem = s.elements[idx]
	prevElem = s.PrevElem(elem)
	conductivity = elem.Material.Conductivity
	deltaX = elem.DistanceFromElement(prevElem)

	aValue = 2 * conductivity / deltaX
	aNeighbour = conductivity / deltaX
	A.SetBand(idx, idx-1, -aNeighbour)
	A.SetBand(idx, idx, aValue)
}

// GetSourceVector return a VecDense used afterward in matrix inversion
func (s *Simulation) SetSourceVector() {
	B := s.solver1D.B
	for _, elem := range s.elements {
		B.SetVec(elem.eNum, elem.Source.Constant*elem.length)
	}
}

// SetBoundaryCondition setup A matric and B vector according to the local BC
func (s *Simulation) SetBoundaryConditions() {
	nElem := len(s.elements)
	A := s.solver1D.A
	B := s.solver1D.B

	for _, bC := range s.boundaryConditions {
		// First Element
		idx := 0
		firstElem := s.elements[0]
		if firstElem.xPosition >= bC.Start && firstElem.xPosition <= bC.End {
			switch bC.BCType {
			case "constant_temperature":
				A.SetBand(idx, idx, 1)
				A.SetBand(idx, idx+1, 0)
				B.SetVec(idx, bC.Value)

			case "constant_heat_flux":
				conductivity := firstElem.Material.Conductivity
				deltaX := s.NextElem(firstElem).DistanceFromElement(firstElem)
				aValue := (conductivity / deltaX) - firstElem.Source.Linear*firstElem.length
				aNeighbour := conductivity / deltaX
				A.SetBand(idx, idx, aValue)
				A.SetBand(idx, idx+1, -aNeighbour)
				B.SetVec(idx, bC.Value+firstElem.Source.Constant*firstElem.length)

			case "convection":
				conductivity := firstElem.Material.Conductivity
				deltaX := s.NextElem(firstElem).DistanceFromElement(firstElem)
				htc := bC.Value
				tFluid := bC.Value2
				aValue := (conductivity / deltaX) + htc - firstElem.Source.Linear*firstElem.length
				aNeighbour := conductivity / deltaX
				A.SetBand(idx, idx, aValue)
				A.SetBand(idx, idx+1, -aNeighbour)
				B.SetVec(idx, htc*tFluid+firstElem.Source.Constant*firstElem.length)
			}
		}

		for idx = 1; idx < nElem-1; idx++ {
			elem := s.elements[idx]
			if elem.xPosition >= bC.Start && elem.xPosition <= bC.End {
				switch bC.BCType {
				case "constant_temperature":
					A.SetBand(idx, idx-1, 0)
					A.SetBand(idx, idx, 1)
					A.SetBand(idx, idx+1, 0)
					B.SetVec(idx, bC.Value)

				case "constant_heat_flux":
					conductivity := elem.Material.Conductivity
					deltaX := s.NextElem(elem).DistanceFromElement(elem)
					aValue := 2*(conductivity/deltaX) - elem.Source.Linear*elem.length
					aNeighbour := conductivity / deltaX
					A.SetBand(idx, idx-1, -aNeighbour)
					A.SetBand(idx, idx, aValue)
					A.SetBand(idx, idx+1, -aNeighbour)
					B.SetVec(idx, bC.Value+elem.Source.Constant*elem.length)

				case "convection":
					conductivity := elem.Material.Conductivity
					deltaX := s.NextElem(elem).DistanceFromElement(elem)
					htc := bC.Value
					tFluid := bC.Value2
					aValue := (conductivity / deltaX) + htc - elem.Source.Linear*elem.length
					aNeighbour := conductivity / deltaX
					A.SetBand(idx, idx-1, -aNeighbour)
					A.SetBand(idx, idx, aValue)
					A.SetBand(idx, idx+1, -aNeighbour)
					B.SetVec(idx, htc*tFluid+elem.Source.Constant*elem.length)
				}
			}
		}

		// Last Element
		lastElem := s.elements[idx]
		if lastElem.xPosition >= bC.Start && lastElem.xPosition <= bC.End+1e-9 {
			switch bC.BCType {
			case "constant_temperature":
				A.SetBand(idx, idx-1, 0)
				A.SetBand(idx, idx, 1)
				B.SetVec(idx, bC.Value)

			case "constant_heat_flux":
				conductivity := lastElem.Material.Conductivity
				deltaX := lastElem.length
				aValue := (conductivity / deltaX) - lastElem.Source.Linear*lastElem.length
				aNeighbour := conductivity / deltaX
				A.SetBand(idx, idx-1, -aNeighbour)
				A.SetBand(idx, idx, aValue)
				B.SetVec(idx, bC.Value+lastElem.Source.Constant*lastElem.length)

			case "convection":
				conductivity := lastElem.Material.Conductivity
				deltaX := s.NextElem(lastElem).DistanceFromElement(lastElem)
				htc := bC.Value
				tFluid := bC.Value2
				aValue := (conductivity / deltaX) + htc - lastElem.Source.Linear*lastElem.length
				aNeighbour := conductivity / deltaX
				A.SetBand(idx, idx-1, -aNeighbour)
				A.SetBand(idx, idx, aValue)
				B.SetVec(idx, htc*tFluid+lastElem.Source.Constant*lastElem.length)
			}
		}
	}
}
