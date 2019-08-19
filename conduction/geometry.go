package conduction

import "math"

type Geometry struct {
	Diameter float64 `json:"diameter"`
	//CrossSectionalArea float64 // [m²] Element cross sectional area
}

// CrossSectionalArea returns the cross sectional area in [m²]
func (g *Geometry) CrossSectionalArea() float64 {
	return math.Pi * math.Pow(g.Diameter, 2) / 4
}

func (g *Geometry) Perimeter() float64 {
	return math.Pi * g.Diameter
}
