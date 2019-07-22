package main

import (
	"conduction/conduction"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"log"
)

func main() {
	caseFileName := "case3_go.json"

	// Load simulation file
	simInput := conduction.NewSimulationInput(caseFileName)

	simulation := conduction.New1DConductionSimulation(simInput)
	simulation.Start()

	simulation.ExportData(caseFileName)

	p, err := plot.New()
	if err != nil {
		log.Panic(err)
	}

	xys := make(plotter.XYs, simInput.Parameters.NElement)
	for i := range xys {
		xys[i].X = simulation.PositionAt(i)
		xys[i].Y = simulation.TemperatureAt(i)
	}

	p.Title.Text = "1D Heat Conduction Steady-State"
	p.X.Label.Text = "X Position"
	p.Y.Label.Text = "Temperature [°C]"

	err = plotutil.AddLinePoints(p, xys)
	if err != nil {
		panic(err)
	}

	// Save the plot to a PNG file.
	if err := p.Save(4*vg.Inch, 4*vg.Inch, "points.png"); err != nil {
		panic(err)
	}
}

// Finite Volume implementation of a 1D steady heat conduction problem
//func main() {
//	caseFileName := "config.json"
//
//	simCase := conduction.LoadCase(caseFileName)
//
//	conductivity := simCase.Geometry.Material.Conductivity
//	deltaX := simCase.NumericalParameters.DeltaX
//	iterationMax := simCase.NumericalParameters.IterationMax
//	tolerance := simCase.NumericalParameters.Tolerance
//
//	gridX := conduction.NewGrid(simCase)
//	nX := len(gridX)
//
//	// Initialize matrix
//	tmpVec := make([]float64, nX)
//	tmpMat := mat.NewDense(nX, nX, nil)
//	A := mat.NewDense(nX, nX, nil)
//	B := mat.NewVecDense(nX, nil)
//	T := mat.NewVecDense(nX, nil)
//
//	// Initialize A matrix
//	for i := 0; i < nX; i++ {
//		A.Set(i, i, 2*conductivity/deltaX)
//		if i > 0 {
//			A.Set(i, i-1, -conductivity/deltaX)
//		}
//		if i < nX-1 {
//			A.Set(i, i+1, -conductivity/deltaX)
//		}
//	}
//
//	// Apply boundary conditions (and override the matrix for Dirichlet conditions)
//	for _, bC := range simCase.BoundaryConditions {
//		 if bC.BCType == "constant_heat_flux" {
//			indexes := conduction.IndexesOf(bC.XPosition, deltaX)
//			if len(indexes) > 1 {
//				for i := indexes[0]; i <= indexes[1]; i++ {
//					fmt.Printf("Idx: %v \n", i)
//					prevValue := B.AtVec(i)
//					B.SetVec(i, prevValue + bC.Value)
//				}
//			} else {
//				i := indexes[0]
//				prevValue := B.AtVec(i)
//				B.SetVec(i, prevValue + bC.Value)
//			}
//		}
//	}
//
//	for _, bC := range simCase.BoundaryConditions {
//		if bC.BCType == "constant_temperature" {
//			indexes := conduction.IndexesOf(bC.XPosition, deltaX)
//			for _, i := range indexes {
//				A.Set(i, i, 1)
//				B.SetVec(i, bC.Value)
//				if i > 0 {
//					A.Set(i, i-1, 0)
//				}
//				if i < nX-1 {
//					A.Set(i, i+1, 0)
//				}
//			}
//		}
//	}
//
//
//	// Solution initialization
//	prevT := mat.NewVecDense(nX, nil)
//	deltaT := 1000.0
//	deltaTVector := make([]float64, iterationMax)
//	iteration := 0
//
//	for deltaT > tolerance && iteration < iterationMax {
//
//
//		PrintDense(A)
//		fmt.Printf("B: %v", B)
//		err := tmpMat.Inverse(A)
//		if err != nil {
//			panic(err)
//		}
//
//		T.MulVec(tmpMat, B)
//		deltaT = floats.Norm(floats.SubTo(tmpVec, T.RawVector().Data, prevT.RawVector().Data), 2)
//		prevT = T
//
//		iteration += 1
//
//	}
//
//	deltaTVector = append(deltaTVector, deltaT)
//	fmt.Printf("Iteration: %v, DeltaT: %v degC \n ", iteration, deltaT)
//	fmt.Printf("T Vector: %v \n ", T)
//
//	p, err := plot.New()
//	if err != nil {
//		log.Panic(err)
//	}
//
//	dataT := T.RawVector().Data
//	xys := make(plotter.XYs, nX)
//	for i := range xys {
//		xys[i].X = gridX[i]
//		xys[i].Y = dataT[i]
//	}
//
//	p.Title.Text = "1D Heat Conduction Steady-State"
//	p.X.Label.Text = "X Position"
//	p.Y.Label.Text = "Temperature [°C]"
//
//	err = plotutil.AddLinePoints(p, xys)
//	if err != nil {
//		panic(err)
//	}
//
//	// Save the plot to a PNG file.
//	if err := p.Save(4*vg.Inch, 4*vg.Inch, "points.png"); err != nil {
//		panic(err)
//	}
//}
