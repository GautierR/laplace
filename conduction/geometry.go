package conduction

type Domain struct {
	Start    float64    `json:"start"`
	End      float64    `json:"end"`
	Source   HeatSource `json:"source"`
	Material Material   `json:"material"`
}

type Material struct {
	Name         string  `json:"name"`
	Density      float64 `json:"density"`
	Conductivity float64 `json:"conductivity"`
	SpecificHeat float64 `json:"specific_heat"`
}
