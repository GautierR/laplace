package conduction

type BoundaryCondition struct {
	Start float64 `json:"start"`
	End   float64 `json:"end"`
	Tag   string  `json:"tag"`
}

func (bC *BoundaryCondition) Length() float64 {
	return bC.End - bC.Start
}

func (bC *BoundaryCondition) IsAtPoint() bool {
	return bC.Start == bC.End
}
