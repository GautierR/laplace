package conduction

import "gonum.org/v1/gonum/mat"

func (e *Element) UpdateInterfaceConductivity(interfaceMethod string) {
	var interpolationFactorW float64 // [-] Interpolation factor at W: (deltaW- / deltaW)
	var interpolationFactorE float64 // [-] Interpolation factor at E: (deltaE+ / deltaE)
	var deltaWMinus float64          // [m] Distance between w (interface) and P
	var deltaEPlus float64           // [m] Distance between P and e
	var condW float64                // [W/m.K] West element thermal conductivity
	var condE float64                // [W/m.K] East element thermal conductivity
	var condP float64                // [W/m.K] Current element thermal conductivity
	condP = e.Material.Conductivity

	if e.previousElement != nil {
		e.deltaW = e.DistanceFromElement(e.previousElement)
		deltaWMinus = e.deltaW - (e.xPosition - e.start)
		interpolationFactorW = deltaWMinus / e.deltaW
		condW = e.previousElement.Material.Conductivity
		e.interfaceConductivityW = InterfaceConductivity(interpolationFactorW, condP, condW, interfaceMethod)
	}

	if e.nextElement != nil {
		e.deltaE = e.DistanceFromElement(e.nextElement)
		deltaEPlus = e.deltaE - (e.end - e.xPosition)
		interpolationFactorE = deltaEPlus / e.deltaE
		condE = e.nextElement.Material.Conductivity
		e.interfaceConductivityE = InterfaceConductivity(interpolationFactorE, condP, condE, interfaceMethod)
	}
}

func InterfaceConductivity(interpolationFactor, conductivity,
	neighborConductivity float64, method string) float64 {

	var interfaceConductivity float64
	switch method {
	case "arithmetic_mean":
		interfaceConductivity = interpolationFactor*conductivity +
			(1-interpolationFactor)*neighborConductivity

	case "harmonic_mean":
		interfaceConductivity = 1 / ((1-interpolationFactor)/conductivity +
			interpolationFactor/neighborConductivity)
	}
	return interfaceConductivity
}

func (e *Element) LinkToMatrices(A *mat.BandDense, B *mat.VecDense,
	state *mat.VecDense, prevState *mat.VecDense) {
	nElement := len(B.RawVector().Data)
	index := e.eNum
	matrixIndex := 3*index + 1

	if index != 0 {
		e.aW = &A.RawBand().Data[matrixIndex-1]
	}

	e.aP = &A.RawBand().Data[matrixIndex]
	e.b = &B.RawVector().Data[index]

	e.State.temperature = &state.RawVector().Data[index]
	e.prevState.temperature = &prevState.RawVector().Data[index]

	if index != nElement-1 {
		e.aE = &A.RawBand().Data[matrixIndex+1]
	}
}
