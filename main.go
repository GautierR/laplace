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
	caseFileName := "case2_go.json"

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
